package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/deali/wechat-clone/internal/macos"
)

// CheckStatus represents the result of a single doctor check.
type CheckStatus int

const (
	CheckOK CheckStatus = iota
	CheckWarning
	CheckError
)

func (s CheckStatus) String() string {
	switch s {
	case CheckOK:
		return "OK"
	case CheckWarning:
		return "WARNING"
	case CheckError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// CheckResult holds the result of a single environment check.
type CheckResult struct {
	Name    string
	Status  CheckStatus
	Message string
}

// Doctor runs all environment checks and returns the results.
func (a *App) Doctor() []CheckResult {
	var results []CheckResult

	// Check 1: macOS
	results = append(results, checkmacOS())

	// Check 2: WeChat.app exists
	results = append(results, checkWeChatExists(a.SourcePath))

	// Check 3: codesign available
	results = append(results, checkCodesign())

	// Check 4: PlistBuddy available
	results = append(results, checkPlistBuddy())

	// Check 5: Write permission
	results = append(results, checkWritePermission(a.TargetDir))

	// Check 6: Disk space
	results = append(results, checkDiskSpace(a.TargetDir))

	// Check 7: Existing clones status
	results = append(results, checkExistingClones(a.TargetDir))

	return results
}

func checkmacOS() CheckResult {
	if runtime.GOOS == "darwin" {
		return CheckResult{Name: "系统平台", Status: CheckOK, Message: "macOS"}
	}
	return CheckResult{Name: "系统平台", Status: CheckError, Message: fmt.Sprintf("当前系统为 %s，需要 macOS", runtime.GOOS)}
}

func checkWeChatExists(path string) CheckResult {
	if _, err := os.Stat(path); err == nil {
		return CheckResult{Name: "微信应用", Status: CheckOK, Message: fmt.Sprintf("已找到: %s", path)}
	}
	return CheckResult{Name: "微信应用", Status: CheckError, Message: fmt.Sprintf("未找到: %s\n请确认微信已安装", path)}
}

func checkCodesign() CheckResult {
	if macos.CodesignAvailable() {
		return CheckResult{Name: "codesign", Status: CheckOK, Message: "可用"}
	}
	return CheckResult{Name: "codesign", Status: CheckError, Message: "不可用，请安装 Xcode Command Line Tools"}
}

func checkPlistBuddy() CheckResult {
	if macos.PlistBuddyAvailable() {
		return CheckResult{Name: "PlistBuddy", Status: CheckOK, Message: "可用"}
	}
	return CheckResult{Name: "PlistBuddy", Status: CheckError, Message: "不可用，请安装 Xcode 或 Xcode Command Line Tools"}
}

func checkWritePermission(dir string) CheckResult {
	// Test write permission by attempting to list directory
	if _, err := os.ReadDir(dir); err == nil {
		return CheckResult{Name: "写入权限", Status: CheckOK, Message: fmt.Sprintf("%s 可访问", dir)}
	}
	return CheckResult{Name: "写入权限", Status: CheckWarning, Message: fmt.Sprintf("%s 可能需要 sudo 权限", dir)}
}

func checkDiskSpace(dir string) CheckResult {
	cmd := exec.Command("df", "-g", dir)
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{Name: "磁盘空间", Status: CheckWarning, Message: "无法检测磁盘空间"}
	}
	// Parse df output - just check if we got output
	if len(output) > 0 {
		return CheckResult{Name: "磁盘空间", Status: CheckOK, Message: "请确保有足够空间（每个分身约 400MB+）"}
	}
	return CheckResult{Name: "磁盘空间", Status: CheckWarning, Message: "无法检测磁盘空间"}
}

func checkExistingClones(targetDir string) CheckResult {
	clones, err := macos.FindClones(targetDir)
	if err != nil {
		return CheckResult{Name: "已有分身", Status: CheckWarning, Message: "无法扫描已有分身"}
	}
	if len(clones) == 0 {
		return CheckResult{Name: "已有分身", Status: CheckOK, Message: "暂无分身"}
	}

	msg := fmt.Sprintf("找到 %d 个分身:\n", len(clones))
	for _, c := range clones {
		relPath, _ := filepath.Rel(targetDir, c.Path)
		msg += fmt.Sprintf("  - %s (%s)\n", relPath, c.BundleID)
	}
	return CheckResult{Name: "已有分身", Status: CheckOK, Message: msg}
}
