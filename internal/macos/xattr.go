package macos

import (
	"fmt"
	"os/exec"
)

// RemoveQuarantine removes the com.apple.quarantine extended attribute from the app.
func RemoveQuarantine(appPath string) error {
	cmd := exec.Command("xattr", "-rd", "com.apple.quarantine", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("清除 quarantine 失败: %s\n%s", err, string(output))
	}
	return nil
}
