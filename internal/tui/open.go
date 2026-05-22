package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/macos"
)

type openState int

const (
	openStateSelect openState = iota
	openStateGuide
)

type openLoadedMsg struct {
	clones []macos.CloneInfo
	err    error
}

type openRevealMsg struct {
	err error
}

type OpenModel struct {
	app    *app.App
	state  openState
	clones []macos.CloneInfo
	cursor int
	all    bool
	err    error
	done   bool
}

func NewOpenModel(a *app.App) *OpenModel {
	return &OpenModel{
		app:   a,
		state: openStateSelect,
		all:   true,
	}
}

func (m *OpenModel) Init() tea.Cmd {
	return func() tea.Msg {
		clones, err := m.app.ListClones()
		return openLoadedMsg{clones: clones, err: err}
	}
}

func (m *OpenModel) Done() bool {
	return m.done
}

func (m *OpenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case openStateSelect:
			switch msg.String() {
			case "up", "k":
				if !m.all && m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if !m.all && m.cursor < len(m.clones)-1 {
					m.cursor++
				}
			case "a":
				m.all = !m.all
			case "enter":
				m.state = openStateGuide
				return m, nil
			}
		case openStateGuide:
			switch msg.String() {
			case "r":
				// Reveal selected/all in Finder
				if m.all {
					for _, c := range m.clones {
						macos.RevealInFinder(c.Path)
					}
				} else if len(m.clones) > 0 {
					macos.RevealInFinder(m.clones[m.cursor].Path)
				}
				return m, nil
			case "enter", "esc":
				m.done = true
				return m, nil
			}
		}

	case openLoadedMsg:
		m.clones = msg.clones
		m.err = msg.err
	}
	return m, nil
}

func (m *OpenModel) View() string {
	title := titleStyle.Render("启动微信分身")

	if m.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			errorStyle.Render(fmt.Sprintf("错误: %v", m.err)),
			"",
			statusStyle.Render("按 Enter 返回主菜单"),
		)
	}

	switch m.state {
	case openStateSelect:
		if len(m.clones) == 0 {
			return lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				infoStyle.Render("暂无分身。请先创建分身。"),
				"",
				statusStyle.Render("按 Enter 返回主菜单"),
			)
		}

		items := ""
		cursor := "  "
		if m.all {
			cursor = "▸ "
		}
		items += menuItemActiveStyle.Render(cursor+"查看所有分身") + "\n"

		for i, c := range m.clones {
			cursor := "  "
			if !m.all && m.cursor == i {
				cursor = "▸ "
			}
			style := menuItemStyle
			if !m.all && m.cursor == i {
				style = menuItemActiveStyle
			}
			items += style.Render(cursor+c.Name) + "\n"
		}

		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			"选择要查看的分身:",
			"",
			items,
			"",
			statusStyle.Render("↑/↓ 选择  •  a 切换全部  •  Enter 确认  •  Esc 返回"),
		)

	case openStateGuide:
		content := title + "\n\n"
		content += "请在 Finder 或启动台中手动打开以下微信分身:\n\n"

		if m.all {
			for _, c := range m.clones {
				content += "  " + lipgloss.NewStyle().Foreground(primary).Render(c.Name) + "\n"
				content += "  " + infoStyle.Render(c.Path) + "\n\n"
			}
		} else if len(m.clones) > 0 {
			c := m.clones[m.cursor]
			content += "  " + lipgloss.NewStyle().Foreground(primary).Render(c.Name) + "\n"
			content += "  " + infoStyle.Render(c.Path) + "\n\n"
		}

		content += infoStyle.Render("提示: 双击 .app 即可启动，也可以拖到程序坞固定") + "\n\n"
		content += statusStyle.Render("r 在 Finder 中定位  •  Enter 返回主菜单")
		return content
	}

	return ""
}
