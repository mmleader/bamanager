package manager

import (
	"bsmanager/backend/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Manager struct {
	Instances map[string]*models.BrowserInstance
	mu        sync.RWMutex
	dataPath  string
	initError error
}

func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".bsmanager")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	m := &Manager{
		Instances: make(map[string]*models.BrowserInstance),
		dataPath:  filepath.Join(configDir, "instances.json"),
	}
	if err := m.load(); err != nil {
		fmt.Printf("Error loading instances: %v\n", err)
		m.initError = err
	}
	go m.startStatusWatcher()
	return m
}

func (m *Manager) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read file error: %w", err)
	}
	if err := json.Unmarshal(data, &m.Instances); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	// Migration: Check if SortNum is 0 and Name has "ID.Name" pattern
	migrated := false
	for _, inst := range m.Instances {
		if inst.SortNum == 0 && strings.Contains(inst.Name, ".") {
			parts := strings.SplitN(inst.Name, ".", 2)
			if len(parts) == 2 {
				if num, err := strconv.Atoi(parts[0]); err == nil {
					inst.SortNum = num
					inst.Name = parts[1] // Keep only the name part
					migrated = true
				}
			}
		}
	}

	if migrated {
		fmt.Println("Migrated instances to use SortNum, saving...")
		// We can't call save() directly because we are holding the lock and save() might not need lock but it's called inside load which holds lock.
		// Actually save() doesn't lock but load() holds lock. safe to call save() ?
		// m.save() implementation: checks m.initError.
		// Let's implement saveInternal to avoid issues or just write file directly here or risk it.
		// Wait, save() doesn't lock. It's safe to call.
		if err := m.save(); err != nil {
			fmt.Printf("Error saving migrated data: %v\n", err)
		}
	}

	return nil
}

func (m *Manager) save() error {
	if m.initError != nil {
		fmt.Println("Warning: Skipping save due to initialization error to prevent data loss")
		return fmt.Errorf("cannot save: initialization failed: %w", m.initError)
	}
	data, err := json.MarshalIndent(m.Instances, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.dataPath, data, 0644)
}

func (m *Manager) ListInstances() ([]*models.BrowserInstance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.initError != nil {
		return nil, m.initError
	}

	var result []*models.BrowserInstance
	for _, inst := range m.Instances {
		result = append(result, inst)
	}

	// Sort by SortNum
	sort.Slice(result, func(i, j int) bool {
		return result[i].SortNum < result[j].SortNum
	})

	return result, nil
}

func (m *Manager) AddInstance(sortNum int, name, path, userDataDir string, args, tags []string) (*models.BrowserInstance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inst := &models.BrowserInstance{
		ID:          uuid.New().String(),
		SortNum:     sortNum,
		Name:        name,
		Path:        path,
		UserDataDir: userDataDir,
		Args:        args,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.Instances[inst.ID] = inst
	err := m.save()
	return inst, err
}

func (m *Manager) UpdateInstance(inst *models.BrowserInstance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Instances[inst.ID]; !ok {
		return fmt.Errorf("instance not found")
	}
	inst.UpdatedAt = time.Now()
	m.Instances[inst.ID] = inst
	return m.save()
}

func (m *Manager) DeleteInstance(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if inst, ok := m.Instances[id]; ok {
		if inst.Running {
			return fmt.Errorf("cannot delete a running instance")
		}
		delete(m.Instances, id)
		return m.save()
	}
	return fmt.Errorf("instance not found")
}

func (m *Manager) StartInstance(id string) error {
	m.mu.Lock()
	inst, ok := m.Instances[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("instance not found")
	}

	if inst.Running && isProcessRunning(inst.PID) {
		m.mu.Unlock()
		return nil
	}

	args := append([]string{}, inst.Args...)
	if inst.UserDataDir != "" {
		args = append(args, fmt.Sprintf("--user-data-dir=%s", inst.UserDataDir))
	}

	cmd := exec.Command(inst.Path, args...)
	err := cmd.Start()
	if err != nil {
		m.mu.Unlock()
		return err
	}

	inst.PID = cmd.Process.Pid
	inst.Running = true
	m.mu.Unlock()

	// Save state
	m.save()

	go func(id string, pid int) {
		cmd.Wait()
		m.mu.Lock()
		if i, ok := m.Instances[id]; ok && i.PID == pid {
			i.Running = false
			i.PID = 0
			m.save()
		}
		m.mu.Unlock()
	}(id, inst.PID)

	return nil
}

func (m *Manager) StopInstance(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inst, ok := m.Instances[id]
	if !ok {
		return fmt.Errorf("instance not found")
	}

	if !inst.Running {
		return nil
	}

	process, err := os.FindProcess(inst.PID)
	if err != nil {
		inst.Running = false
		return nil
	}

	// 尝试优雅关闭，如果不行再 Kill
	if runtime.GOOS != "windows" {
		process.Signal(syscall.SIGTERM)
		time.Sleep(500 * time.Millisecond)
	}

	process.Kill()
	inst.Running = false
	inst.PID = 0
	m.save()
	return nil
}

func (m *Manager) startStatusWatcher() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		m.mu.Lock()
		changed := false
		for _, inst := range m.Instances {
			if inst.Running {
				if !isProcessRunning(inst.PID) {
					inst.Running = false
					inst.PID = 0
					changed = true
				}
			}
		}
		if changed {
			m.save()
		}
		m.mu.Unlock()
	}
}

func isProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}

	if runtime.GOOS != "windows" {
		process, err := os.FindProcess(pid)
		if err != nil {
			return false
		}
		err = process.Signal(syscall.Signal(0))
		return err == nil
	}

	// Windows 简单处理， FindProcess 在 Windows 上不保证存活
	// 实际开发中可以使用更多 Windows API
	process, err := os.FindProcess(pid)
	return err == nil && process != nil
}

func (m *Manager) CheckInstanceProxy(id, target string) (map[string]interface{}, error) {
	m.mu.Lock()
	inst, ok := m.Instances[id]
	if !ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("instance not found")
	}
	args := inst.Args
	m.mu.Unlock() // Unlock early to perform network request

	// Parse proxy from args
	var proxyURL string
	for _, arg := range args {
		if len(arg) > 15 && arg[:15] == "--proxy-server=" {
			proxyURL = arg[15:]
			break
		}
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if proxyURL != "" {
		if !strings.HasPrefix(proxyURL, "http") && !strings.HasPrefix(proxyURL, "socks") {
			proxyURL = "http://" + proxyURL
		}
		u, err := url.Parse(proxyURL)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
	}

	var checkURL string
	if target == "cn" {
		checkURL = "https://myip.ipip.net"
	} else {
		// global (default)
		checkURL = "http://ip-api.com/json"
	}

	start := time.Now()
	resp, err := client.Get(checkURL)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return nil, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	var region string
	var detail string

	if target == "cn" {
		// myip.ipip.net returns text: "当前 IP：1.2.3.4  来自于：中国 广东 深圳  电信"
		detail = strings.TrimSpace(content)
		parts := strings.Split(detail, "来自于：")
		if len(parts) > 1 {
			locInfo := strings.TrimSpace(parts[1])
			// Split by spaces to get country, province, city
			locParts := strings.Fields(locInfo)
			if len(locParts) > 0 {
				country := locParts[0]
				code := "UNKNOWN"
				switch country {
				case "中国":
					code = "CN"
				case "香港", "中国香港":
					code = "HK"
				case "澳门", "中国澳门":
					code = "MO"
				case "台湾", "中国台湾":
					code = "TW"
				default:
					// If it's something else, try to use first 2 chars of country or mapped manually if needed
					// For now, just keep original first word if not matched
					code = country
				}

				region = code
				// Append province if available and not same as country (e.g. for CN)
				if len(locParts) > 1 && (code == "CN" || code == "TW") {
					region += " " + locParts[1]
				}
			} else {
				region = locInfo
			}
		} else {
			region = detail
		}
	} else {
		// ip-api.com returns JSON
		detail = content
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			if code, ok := result["countryCode"].(string); ok {
				region = code
			} else if country, ok := result["country"].(string); ok {
				region = country
			} else {
				region = "Unknown"
			}
		} else {
			region = "Parse Error"
		}
	}

	// Update instance
	m.mu.Lock()
	if inst, ok := m.Instances[id]; ok {
		inst.ProxyRegion = region
		inst.ProxyLatency = latency
		inst.ProxyDetail = detail
		m.save()
	}
	m.mu.Unlock()

	return map[string]interface{}{
		"region":  region,
		"latency": latency,
		"detail":  detail,
	}, nil
}
