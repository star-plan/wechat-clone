package macos

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetBundleID modifies the CFBundleIdentifier in the given Info.plist using PlistBuddy.
func SetBundleID(plistPath, bundleID string) error {
	cmd := exec.Command("sudo", "/usr/libexec/PlistBuddy",
		"-c", fmt.Sprintf("Set :CFBundleIdentifier %s", bundleID),
		plistPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("修改 Bundle ID 失败: %s\n%s", err, string(output))
	}
	return nil
}

// GetBundleID reads the CFBundleIdentifier from the given Info.plist.
func GetBundleID(plistPath string) (string, error) {
	cmd := exec.Command("/usr/libexec/PlistBuddy",
		"-c", "Print :CFBundleIdentifier",
		plistPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("读取 Bundle ID 失败: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// PlistBuddyAvailable checks if PlistBuddy is available on the system.
func PlistBuddyAvailable() bool {
	_, err := exec.LookPath("/usr/libexec/PlistBuddy")
	return err == nil
}
