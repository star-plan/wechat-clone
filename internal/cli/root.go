package cli

import (
	"fmt"
	"os"

	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var appInstance *app.App

var rootCmd = &cobra.Command{
	Use:   "wechat-clone",
	Short: "macOS 微信分身管理工具",
	Long:  "wechat-clone 是一个 macOS 微信分身管理 CLI 工具，支持创建、管理、启动多个微信分身。",
	Run: func(cmd *cobra.Command, args []string) {
		// Default: launch TUI
		m := tui.NewMainModel(appInstance)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "TUI 启动失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	appInstance = app.New()
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
