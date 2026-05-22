package app

import (
	"fmt"

	"github.com/deali/wechat-clone/internal/macos"
)

// OpenClones starts the specified clones. If no IDs are given, all clones are started.
func (a *App) OpenClones(ids ...int) error {
	clones, err := a.ListClones()
	if err != nil {
		return fmt.Errorf("获取分身列表失败: %w", err)
	}
	if len(clones) == 0 {
		return fmt.Errorf("没有找到任何分身，请先使用 create 命令创建")
	}

	// If no specific IDs, open all
	if len(ids) == 0 {
		for _, c := range clones {
			if err := macos.LaunchApp(c.Path); err != nil {
				return fmt.Errorf("启动 %s 失败: %w", c.Name, err)
			}
		}
		return nil
	}

	// Open specific clones
	for _, id := range ids {
		found := false
		for _, c := range clones {
			if c.ID == id {
				if err := macos.LaunchApp(c.Path); err != nil {
					return fmt.Errorf("启动 %s 失败: %w", c.Name, err)
				}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("未找到编号为 %d 的分身", id)
		}
	}

	return nil
}
