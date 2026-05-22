package macos

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// LaunchApp starts a WeChat clone. It tries nohup first, then falls back to open.
func LaunchApp(appPath string) error {
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return fmt.Errorf("应用不存在: %s", appPath)
	}

	// Try nohup first (proven approach from community)
	binaryPath := filepath.Join(appPath, "Contents", "MacOS", "WeChat")
	if _, err := os.Stat(binaryPath); err == nil {
		cmd := exec.Command("nohup", binaryPath)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
		if err := cmd.Start(); err == nil {
			return nil
		}
	}

	// Fallback to open command
	cmd := exec.Command("open", appPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动应用失败: %w", err)
	}
	return nil
}
