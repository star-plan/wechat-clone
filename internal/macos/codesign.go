package macos

import (
	"fmt"
	"os/exec"
)

// ResignApp performs an ad-hoc codesign on the given app bundle.
func ResignApp(appPath string) error {
	cmd := exec.Command("codesign", "--force", "--deep", "--sign", "-", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("签名失败: %s\n%s", err, string(output))
	}
	return nil
}

// CodesignAvailable checks if codesign is available on the system.
func CodesignAvailable() bool {
	_, err := exec.LookPath("codesign")
	return err == nil
}
