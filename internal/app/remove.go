package app

import (
	"fmt"

	"github.com/deali/wechat-clone/internal/macos"
)

// RemoveResult holds the result of removing a single clone.
type RemoveResult struct {
	ID     int
	Path   string
	Status string // "removed", "error"
	Err    error
}

// RemoveClones removes the specified clones. If ids contains -1, all clones are removed.
func (a *App) RemoveClones(ids []int, force bool) ([]RemoveResult, error) {
	clones, err := a.ListClones()
	if err != nil {
		return nil, fmt.Errorf("获取分身列表失败: %w", err)
	}
	if len(clones) == 0 {
		return nil, fmt.Errorf("没有找到任何分身")
	}

	// Determine targets
	var targets []macos.CloneInfo
	removeAll := false
	for _, id := range ids {
		if id == -1 {
			removeAll = true
			break
		}
	}

	if removeAll {
		targets = clones
	} else {
		for _, id := range ids {
			found := false
			for _, c := range clones {
				if c.ID == id {
					targets = append(targets, c)
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("未找到编号为 %d 的分身", id)
			}
		}
	}

	results := make([]RemoveResult, 0, len(targets))
	for _, c := range targets {
		if err := macos.RemoveApp(c.Path); err != nil {
			results = append(results, RemoveResult{ID: c.ID, Path: c.Path, Status: "error", Err: err})
		} else {
			results = append(results, RemoveResult{ID: c.ID, Path: c.Path, Status: "removed"})
		}
	}

	return results, nil
}
