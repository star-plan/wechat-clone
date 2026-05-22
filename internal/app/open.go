package app

import (
	"fmt"

	"github.com/deali/wechat-clone/internal/macos"
)

// OpenGuide returns the paths of clones for manual launching guidance.
// If no IDs are given, all clones are returned.
func (a *App) OpenGuide(ids ...int) ([]macos.CloneInfo, error) {
	clones, err := a.ListClones()
	if err != nil {
		return nil, fmt.Errorf("获取分身列表失败: %w", err)
	}
	if len(clones) == 0 {
		return nil, fmt.Errorf("没有找到任何分身，请先使用 create 命令创建")
	}

	if len(ids) == 0 {
		return clones, nil
	}

	var result []macos.CloneInfo
	for _, id := range ids {
		found := false
		for _, c := range clones {
			if c.ID == id {
				result = append(result, c)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("未找到编号为 %d 的分身", id)
		}
	}
	return result, nil
}
