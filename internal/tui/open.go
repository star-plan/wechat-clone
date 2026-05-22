package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/macos"
)

type openState int

const (
	openStateSelect openState = iota
	openStateWorking
	openStateDone
)

type openLoadedMsg struct {
	clones []macos.CloneInfo
	err    error
}

type openDoneMsg struct {
	err error
}

type OpenModel struct {
	app     *app.App
	state   openState
	clones  []macos.CloneInfo
	cursor  int
	all     bool
	spinner spinner.Model
	err     error
	done    bool
}

func NewOpenModel(a *app.App) *OpenModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return &OpenModel{
		app:     a,
		state:   openStateSelect,
		all:     true,
		spinner: s,
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
				m.state = openStateWorking
				return m, tea.Batch(m.spinner.Tick, m.doOpen())
			}
		case openStateDone:
			switch msg.String() {
			case "enter":
				m.done = true
				return m, nil
			}
		}

	case openLoadedMsg:
		m.clones = msg.clones
		m.err = msg.err

	case openDoneMsg:
		m.err = msg.err
		m.state = openStateDone
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *OpenModel) doOpen() tea.Cmd {
	return func() tea.Msg {
		if m.all {
			return openDoneMsg{err: m.app.OpenClones()}
		}
		if len(m.clones) == 0 {
			return openDoneMsg{err: fmt.Errorf("没有找到任何分身")}
		}
		id := m.clones[m.cursor].ID
		return openDoneMsg{err: m.app.OpenClones(id)}
	}
}

func (m *OpenModel) View() string {
	title := titleStyle.Render("启动微信分身")

	if m.err != nil && m.state != openStateDone {
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
		// "All" option
		cursor := "  "
		if m.all {
			cursor = "▸ "
		}
		items += menuItemActiveStyle.Render(cursor+"打开所有分身") + "\n"

		// Individual clones
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
			"选择要启动的分身:",
			"",
			items,
			"",
			statusStyle.Render("↑/↓ 选择  •  a 切换全部  •  Enter 确认  •  Esc 返回"),
		)

	case openStateWorking:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " 正在启动..."),
		)

	case openStateDone:
		content := title + "\n\n"
		if m.err != nil {
			content += errorStyle.Render(fmt.Sprintf("错误: %v", m.err))
		} else {
			content += successStyle.Render("启动成功!")
		}
		content += "\n\n" + statusStyle.Render("按 Enter 返回主菜单")
		return content
	}

	return ""
}
