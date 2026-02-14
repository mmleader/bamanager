package models

import "time"

type BrowserInstance struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`        // 浏览器执行程序路径
	UserDataDir string    `json:"userDataDir"` // 用户数据目录（核心指纹隔离）
	Args        []string  `json:"args"`        // 自定义启动参数
	Tags          []string  `json:"tags"`          // 标签（用于分类）
	ProxyRegion   string    `json:"proxyRegion"`   // 代理IP归属地 (简要)
	ProxyLatency  int64     `json:"proxyLatency"`  // 延迟 (ms)
	ProxyDetail   string    `json:"proxyDetail"`   // 详细信息 (Tooltip)
	Running       bool      `json:"running"`       // 是否运行中
	PID         int       `json:"pid"`         // 进程ID
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
