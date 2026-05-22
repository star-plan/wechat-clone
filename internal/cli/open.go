package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/deali/wechat-clone/internal/macos"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [编号|all]",
	Short: "查看微信分身启动指引",
	Long:  "显示微信分身的路径，引导你在 Finder 中手动启动。可选在 Finder 中定位。",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		if len(args) > 0 {
			if args[0] == "all" {
				// Show all
			} else {
				id, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "错误: 无效的编号 '%s'\n", args[0])
					os.Exit(1)
				}
				ids = append(ids, id)
			}
		}

		clones, err := appInstance.OpenGuide(ids...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("请在 Finder 或启动台中手动打开以下微信分身:\n")
		for _, c := range clones {
			fmt.Printf("  %s\n", c.Path)
		}

		fmt.Println("\n提示: 双击 .app 即可启动，也可以拖到程序坞固定。")
		fmt.Println("\n是否在 Finder 中定位？输入编号打开，或直接回车跳过:")

		for i, c := range clones {
			fmt.Printf("  %d) %s\n", i+1, c.Name)
		}
		fmt.Print("\n选择: ")

		var choice int
		if _, err := fmt.Scan(&choice); err == nil && choice >= 1 && choice <= len(clones) {
			if err := macos.RevealInFinder(clones[choice-1].Path); err != nil {
				fmt.Fprintf(os.Stderr, "在 Finder 中定位失败: %v\n", err)
			} else {
				fmt.Println("已在 Finder 中打开。")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
