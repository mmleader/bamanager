package main

import (
	"bsmanager/backend/config"
	"bsmanager/backend/manager"
	"bsmanager/backend/models"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	manager *manager.Manager
	Config  *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		manager: manager.NewManager(),
		Config:  config.NewConfig(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	width, height := runtime.WindowGetSize(ctx)
	a.Config.Width = width
	a.Config.Height = height
	a.Config.Save()
	return false
}

func (a *App) ListInstances() []*models.BrowserInstance {
	return a.manager.ListInstances()
}

func (a *App) AddInstance(name, path, userDataDir string, args []string) *models.BrowserInstance {
	return a.manager.AddInstance(name, path, userDataDir, args)
}

func (a *App) UpdateInstance(inst *models.BrowserInstance) error {
	return a.manager.UpdateInstance(inst)
}

func (a *App) DeleteInstance(id string) error {
	return a.manager.DeleteInstance(id)
}

func (a *App) StartInstance(id string) error {
	return a.manager.StartInstance(id)
}

func (a *App) StopInstance(id string) error {
	return a.manager.StopInstance(id)
}

func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择浏览器执行文件",
	})
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择用户数据目录",
	})
}
