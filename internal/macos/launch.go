package macos

import (
	"fmt"
	"os"
	"os/exec"
)

// RevealInFinder opens Finder and selects the given app path.
func RevealInFinder(appPath string) error {
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return fmt.Errorf("应用不存在: %s", appPath)
	}
	cmd := exec.Command("open", "-R", appPath)
	return cmd.Run()
}
