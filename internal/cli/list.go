package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出当前已有分身",
	Long:  "列出当前已有的所有微信分身及其 Bundle ID。",
	Run: func(cmd *cobra.Command, args []string) {
		clones, err := appInstance.ListClones()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		if len(clones) == 0 {
			fmt.Println("暂无分身。使用 'wechat-clone create <数量>' 创建分身。")
			return
		}

		fmt.Printf("找到 %d 个微信分身:\n\n", len(clones))
		for _, c := range clones {
			fmt.Printf("  %s\t%s\n", c.Name, c.BundleID)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
