package cli

import (
	"fmt"
	"os"

	"github.com/deali/wechat-clone/internal/app"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "检查环境配置",
	Long:  "检查当前系统环境是否满足运行微信分身的所有条件。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("正在检查环境...\n")

		results := appInstance.Doctor()

		allOK := true
		for _, r := range results {
			icon := "✓"
			switch r.Status {
			case app.CheckOK:
				icon = "✓"
			case app.CheckWarning:
				icon = "⚠"
				allOK = false
			case app.CheckError:
				icon = "✗"
				allOK = false
			}
			fmt.Printf("  %s %s: %s\n", icon, r.Name, r.Message)
		}

		fmt.Println()
		if allOK {
			fmt.Println("环境检查通过! 可以正常使用。")
		} else {
			fmt.Println("存在以上问题，请修复后再使用。")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
