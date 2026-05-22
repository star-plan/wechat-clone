package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [编号|all]",
	Short: "打开微信分身",
	Long:  "打开指定编号或所有微信分身。例如: wechat-clone open 2 或 wechat-clone open",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		if len(args) > 0 {
			if args[0] == "all" {
				// Open all - no IDs specified
			} else {
				id, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "错误: 无效的编号 '%s'\n", args[0])
					os.Exit(1)
				}
				ids = append(ids, id)
			}
		}

		if len(ids) == 0 {
			fmt.Println("正在打开所有微信分身...")
		} else {
			fmt.Printf("正在打开微信分身 %d...\n", ids[0])
		}

		if err := appInstance.OpenClones(ids...); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("完成!")
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
