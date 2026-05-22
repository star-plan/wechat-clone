package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var removeForce bool

var removeCmd = &cobra.Command{
	Use:   "remove <编号|all>",
	Short: "删除微信分身",
	Long:  "删除指定编号或所有微信分身。删除前需要确认，除非使用 --force。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		if args[0] == "all" {
			ids = append(ids, -1) // -1 means all
		} else {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: 无效的编号 '%s'\n", args[0])
				os.Exit(1)
			}
			ids = append(ids, id)
		}

		// Confirm unless --force
		if !removeForce {
			fmt.Print("确定要删除分身吗？此操作不可恢复。(y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println("已取消。")
				return
			}
		}

		fmt.Println("正在删除分身...")

		results, err := appInstance.RemoveClones(ids, removeForce)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		for _, r := range results {
			switch r.Status {
			case "removed":
				fmt.Printf("  ✓ 分身 %d 已删除\n", r.ID)
			case "error":
				fmt.Fprintf(os.Stderr, "  ✗ 分身 %d 删除失败: %v\n", r.ID, r.Err)
			}
		}

		fmt.Println("\n完成!")
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "跳过确认直接删除")
	rootCmd.AddCommand(removeCmd)
}
