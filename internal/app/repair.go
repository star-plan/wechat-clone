package app

import (
	"fmt"

	"github.com/deali/wechat-clone/internal/macos"
)

// RepairResult holds the result of repairing a single clone.
type RepairResult struct {
	ID     int
	Path   string
	Status string // "repaired", "error"
	Err    error
}

// RepairClones repairs the specified clones. If no IDs are given, all clones are repaired.
func (a *App) RepairClones(ids ...int) ([]RepairResult, error) {
	clones, err := a.ListClones()
	if err != nil {
		return nil, fmt.Errorf("获取分身列表失败: %w", err)
	}
	if len(clones) == 0 {
		return nil, fmt.Errorf("没有找到任何分身，请先使用 create 命令创建")
	}

	// If no specific IDs, repair all
	targets := clones
	if len(ids) > 0 {
		targets = nil
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

	results := make([]RepairResult, 0, len(targets))
	for _, c := range targets {
		result := a.repairOne(c)
		results = append(results, result)
	}

	return results, nil
}

func (a *App) repairOne(clone macos.CloneInfo) RepairResult {
	// Step 1: Re-set bundle ID
	plistPath := clone.Path + "/Contents/Info.plist"
	bundleID := CloneBundleID(clone.ID)
	if err := macos.SetBundleID(plistPath, bundleID); err != nil {
		return RepairResult{ID: clone.ID, Path: clone.Path, Status: "error", Err: fmt.Errorf("修改 Bundle ID 失败: %w", err)}
	}

	// Step 2: Resign
	if err := macos.ResignApp(clone.Path); err != nil {
		return RepairResult{ID: clone.ID, Path: clone.Path, Status: "error", Err: fmt.Errorf("签名失败: %w", err)}
	}

	// Step 3: Remove quarantine
	if err := macos.RemoveQuarantine(clone.Path); err != nil {
		return RepairResult{ID: clone.ID, Path: clone.Path, Status: "error", Err: fmt.Errorf("清除 quarantine 失败: %w", err)}
	}

	return RepairResult{ID: clone.ID, Path: clone.Path, Status: "repaired"}
}
