package macos

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

const (
	DefaultWeChatPath = "/Applications/WeChat.app"
	DefaultTargetDir  = "/Applications"
	ClonePrefix       = "WeChat Clone"
	CloneSuffix       = ".app"
)

// CloneInfo holds information about a WeChat clone.
type CloneInfo struct {
	ID        int
	Name      string
	Path      string
	BundleID  string
}

// CloneApp copies the source WeChat.app to the destination path using sudo.
func CloneApp(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("源应用不存在: %s", src)
	}
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("目标应用已存在: %s", dst)
	}
	cmd := exec.Command("sudo", "cp", "-R", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveApp removes a clone app directory using sudo.
func RemoveApp(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("应用不存在: %s", path)
	}
	cmd := exec.Command("sudo", "rm", "-rf", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// FindClones scans the target directory for WeChat clones.
func FindClones(targetDir string) ([]CloneInfo, error) {
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("无法读取目录 %s: %w", targetDir, err)
	}

	re := regexp.MustCompile(`^` + regexp.QuoteMeta(ClonePrefix) + ` (\d+)` + regexp.QuoteMeta(CloneSuffix) + `$`)
	var clones []CloneInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matches := re.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}
		appPath := filepath.Join(targetDir, entry.Name())
		plistPath := filepath.Join(appPath, "Contents", "Info.plist")
		bundleID, _ := GetBundleID(plistPath)

		clones = append(clones, CloneInfo{
			ID:       id,
			Name:     entry.Name(),
			Path:     appPath,
			BundleID: bundleID,
		})
	}

	sort.Slice(clones, func(i, j int) bool {
		return clones[i].ID < clones[j].ID
	})

	return clones, nil
}

// NextCloneID returns the next available clone ID (max existing + 1).
func NextCloneID(targetDir string) (int, error) {
	clones, err := FindClones(targetDir)
	if err != nil {
		return 1, err
	}
	if len(clones) == 0 {
		return 1, nil
	}
	return clones[len(clones)-1].ID + 1, nil
}

// ClonePath returns the expected path for a clone with the given ID.
func ClonePath(targetDir string, id int) string {
	name := fmt.Sprintf("%s %d%s", ClonePrefix, id, CloneSuffix)
	return filepath.Join(targetDir, name)
}
