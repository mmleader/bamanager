package manager

import (
	"bsmanager/backend/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ProxyManager struct {
	Proxies map[string]*models.ProxyConfig
	mu      sync.RWMutex
	dataPath string
}

func NewProxyManager() *ProxyManager {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".bsmanager")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	pm := &ProxyManager{
		Proxies:  make(map[string]*models.ProxyConfig),
		dataPath: filepath.Join(configDir, "proxies.json"),
	}
	if err := pm.load(); err != nil {
		fmt.Printf("Error loading proxies: %v\n", err)
	}
	return pm
}

func (pm *ProxyManager) load() error {
	data, err := os.ReadFile(pm.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read proxies file error: %w", err)
	}
	return json.Unmarshal(data, &pm.Proxies)
}

func (pm *ProxyManager) save() error {
	data, err := json.MarshalIndent(pm.Proxies, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pm.dataPath, data, 0644)
}

func (pm *ProxyManager) ListProxies() []*models.ProxyConfig {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []*models.ProxyConfig
	for _, p := range pm.Proxies {
		result = append(result, p)
	}
	return result
}

func (pm *ProxyManager) GetProxy(id string) *models.ProxyConfig {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.Proxies[id]
}

func (pm *ProxyManager) GetProxyURL(id string) string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if p, ok := pm.Proxies[id]; ok {
		return p.ToURL()
	}
	return ""
}

func (pm *ProxyManager) AddProxy(name, protocol, host string, port int, username, password string) (*models.ProxyConfig, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p := &models.ProxyConfig{
		ID:        uuid.New().String(),
		Name:      name,
		Protocol:  protocol,
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	pm.Proxies[p.ID] = p
	err := pm.save()
	return p, err
}

func (pm *ProxyManager) UpdateProxy(proxy *models.ProxyConfig) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.Proxies[proxy.ID]; !ok {
		return fmt.Errorf("proxy not found")
	}
	proxy.UpdatedAt = time.Now()
	pm.Proxies[proxy.ID] = proxy
	return pm.save()
}

func (pm *ProxyManager) DeleteProxy(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.Proxies[id]; !ok {
		return fmt.Errorf("proxy not found")
	}
	delete(pm.Proxies, id)
	return pm.save()
}
