package app

import (
	"fmt"

	"github.com/deali/wechat-clone/internal/macos"
)

// CreateResult holds the result of creating a single clone.
type CreateResult struct {
	ID     int
	Path   string
	Status string // "created", "skipped", "error"
	Err    error
}

// CreateClones creates the specified number of WeChat clones.
func (a *App) CreateClones(count int, force bool) ([]CreateResult, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, fmt.Errorf("创建数量必须大于 0")
	}

	nextID, err := macos.NextCloneID(a.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("获取下一个可用编号失败: %w", err)
	}

	results := make([]CreateResult, 0, count)
	for i := 0; i < count; i++ {
		id := nextID + i
		result := a.createOne(id, force)
		results = append(results, result)
	}

	return results, nil
}

func (a *App) createOne(id int, force bool) CreateResult {
	dst := a.ClonePath(id)

	// Check if already exists
	if _, err := macos.FindClones(a.TargetDir); err == nil {
		clones, _ := macos.FindClones(a.TargetDir)
		for _, c := range clones {
			if c.ID == id {
				if !force {
					return CreateResult{
						ID:     id,
						Path:   dst,
						Status: "skipped",
						Err:    fmt.Errorf("分身已存在，使用 --force 覆盖"),
					}
				}
				// Remove existing clone before recreating
				if err := macos.RemoveApp(dst); err != nil {
					return CreateResult{ID: id, Path: dst, Status: "error", Err: err}
				}
			}
		}
	}

	// Step 1: Copy app
	if err := macos.CloneApp(a.SourcePath, dst); err != nil {
		return CreateResult{ID: id, Path: dst, Status: "error", Err: fmt.Errorf("复制失败: %w", err)}
	}

	// Step 2: Set bundle ID
	plistPath := dst + "/Contents/Info.plist"
	bundleID := CloneBundleID(id)
	if err := macos.SetBundleID(plistPath, bundleID); err != nil {
		return CreateResult{ID: id, Path: dst, Status: "error", Err: fmt.Errorf("修改 Bundle ID 失败: %w", err)}
	}

	// Step 3: Resign
	if err := macos.ResignApp(dst); err != nil {
		return CreateResult{ID: id, Path: dst, Status: "error", Err: fmt.Errorf("签名失败: %w", err)}
	}

	// Step 4: Remove quarantine
	if err := macos.RemoveQuarantine(dst); err != nil {
		return CreateResult{ID: id, Path: dst, Status: "error", Err: fmt.Errorf("清除 quarantine 失败: %w", err)}
	}

	return CreateResult{ID: id, Path: dst, Status: "created"}
}
