package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deali/wechat-clone/internal/macos"
)

// App holds configuration for the wechat-clone tool.
type App struct {
	SourcePath string // Path to original WeChat.app
	TargetDir  string // Directory where clones are placed
}

// New creates a new App with default configuration.
func New() *App {
	return &App{
		SourcePath: macos.DefaultWeChatPath,
		TargetDir:  macos.DefaultTargetDir,
	}
}

// Validate checks if the environment is ready for operations.
func (a *App) Validate() error {
	if _, err := os.Stat(a.SourcePath); os.IsNotExist(err) {
		return fmt.Errorf("未找到微信应用: %s\n请确认微信已安装在 /Applications 目录下", a.SourcePath)
	}
	if _, err := os.Stat(a.TargetDir); os.IsNotExist(err) {
		return fmt.Errorf("目标目录不存在: %s", a.TargetDir)
	}
	return nil
}

// ListClones returns all existing WeChat clones.
func (a *App) ListClones() ([]macos.CloneInfo, error) {
	return macos.FindClones(a.TargetDir)
}

// CloneName returns the display name for a clone with the given ID.
func CloneName(id int) string {
	return fmt.Sprintf("%s %d%s", macos.ClonePrefix, id, macos.CloneSuffix)
}

// CloneBundleID returns the bundle identifier for a clone with the given ID.
func CloneBundleID(id int) string {
	return fmt.Sprintf("com.tencent.xinWeChat.clone%d", id)
}

// ClonePath returns the full path for a clone with the given ID.
func (a *App) ClonePath(id int) string {
	return filepath.Join(a.TargetDir, CloneName(id))
}

// FindCloneByID finds a specific clone by its ID.
func (a *App) FindCloneByID(id int) (*macos.CloneInfo, error) {
	clones, err := a.ListClones()
	if err != nil {
		return nil, err
	}
	for _, c := range clones {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("未找到编号为 %d 的分身", id)
}
