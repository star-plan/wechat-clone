package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var repairCmd = &cobra.Command{
	Use:   "repair [编号|all]",
	Short: "修复微信分身签名",
	Long:  "重新修复指定或所有微信分身的签名和 quarantine 属性。",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		if len(args) > 0 {
			if args[0] == "all" {
				// Repair all - no IDs specified
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
			fmt.Println("正在修复所有微信分身...")
		} else {
			fmt.Printf("正在修复微信分身 %d...\n", ids[0])
		}

		results, err := appInstance.RepairClones(ids...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		for _, r := range results {
			switch r.Status {
			case "repaired":
				fmt.Printf("  ✓ 分身 %d 修复成功\n", r.ID)
			case "error":
				fmt.Fprintf(os.Stderr, "  ✗ 分身 %d 修复失败: %v\n", r.ID, r.Err)
			}
		}

		fmt.Println("\n完成!")
	},
}

func init() {
	rootCmd.AddCommand(repairCmd)
}
