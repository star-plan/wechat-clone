package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var createForce bool

var createCmd = &cobra.Command{
	Use:   "create <数量>",
	Short: "创建指定数量的微信分身",
	Long:  "创建指定数量的微信分身。例如: wechat-clone create 3",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无效的数量 '%s'，请输入数字\n", args[0])
			os.Exit(1)
		}

		fmt.Printf("正在创建 %d 个微信分身...\n\n", count)

		results, err := appInstance.CreateClones(count, createForce)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		for _, r := range results {
			switch r.Status {
			case "created":
				fmt.Printf("  ✓ 分身 %d 创建成功: %s\n", r.ID, r.Path)
			case "skipped":
				fmt.Printf("  - 分身 %d 已存在，跳过: %s\n", r.ID, r.Path)
			case "error":
				fmt.Fprintf(os.Stderr, "  ✗ 分身 %d 创建失败: %v\n", r.ID, r.Err)
			}
		}

		fmt.Println("\n完成!")
	},
}

func init() {
	createCmd.Flags().BoolVarP(&createForce, "force", "f", false, "覆盖已存在的分身")
	rootCmd.AddCommand(createCmd)
}
