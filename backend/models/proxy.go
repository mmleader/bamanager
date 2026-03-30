package models

import "time"

type ProxyConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`     // 代理名称，如 "日本节点"
	Protocol  string    `json:"protocol"` // socks5 / http / https
	Host      string    `json:"host"`     // 代理地址
	Port      int       `json:"port"`     // 端口
	Username  string    `json:"username"` // 认证用户名（可选）
	Password  string    `json:"password"` // 认证密码（可选）
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToURL 返回格式化后的代理URL，如 socks5://user:pass@host:port
func (p *ProxyConfig) ToURL() string {
	auth := ""
	if p.Username != "" {
		auth = p.Username
		if p.Password != "" {
			auth += ":" + p.Password
		}
		auth += "@"
	}
	return p.Protocol + "://" + auth + p.Host + ":" + itoa(p.Port)
}

// itoa converts int to string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
