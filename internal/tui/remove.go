package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/macos"
)

type removeState int

const (
	removeStateSelect removeState = iota
	removeStateConfirm
	removeStateWorking
	removeStateDone
)

type removeLoadedMsg struct {
	clones []macos.CloneInfo
	err    error
}

type removeDoneMsg struct {
	results []app.RemoveResult
	err     error
}

type RemoveModel struct {
	app     *app.App
	state   removeState
	clones  []macos.CloneInfo
	cursor  int
	all     bool
	spinner spinner.Model
	results []app.RemoveResult
	err     error
	done    bool
}

func NewRemoveModel(a *app.App) *RemoveModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return &RemoveModel{
		app:     a,
		state:   removeStateSelect,
		all:     false,
		spinner: s,
	}
}

func (m *RemoveModel) Init() tea.Cmd {
	return func() tea.Msg {
		clones, err := m.app.ListClones()
		return removeLoadedMsg{clones: clones, err: err}
	}
}

func (m *RemoveModel) Done() bool {
	return m.done
}

func (m *RemoveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case removeStateSelect:
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
				m.state = removeStateConfirm
				return m, nil
			}
		case removeStateConfirm:
			switch msg.String() {
			case "y", "Y":
				m.state = removeStateWorking
				return m, tea.Batch(m.spinner.Tick, m.doRemove())
			case "n", "N", "esc":
				m.state = removeStateSelect
				return m, nil
			}
		case removeStateDone:
			switch msg.String() {
			case "enter":
				m.done = true
				return m, nil
			}
		}

	case removeLoadedMsg:
		m.clones = msg.clones
		m.err = msg.err

	case removeDoneMsg:
		m.results = msg.results
		m.err = msg.err
		m.state = removeStateDone
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *RemoveModel) doRemove() tea.Cmd {
	return func() tea.Msg {
		if m.all {
			results, err := m.app.RemoveClones([]int{-1}, true)
			return removeDoneMsg{results: results, err: err}
		}
		if len(m.clones) == 0 {
			return removeDoneMsg{err: fmt.Errorf("没有找到任何分身")}
		}
		id := m.clones[m.cursor].ID
		results, err := m.app.RemoveClones([]int{id}, true)
		return removeDoneMsg{results: results, err: err}
	}
}

func (m *RemoveModel) View() string {
	title := titleStyle.Render("删除微信分身")

	if m.err != nil && m.state != removeStateDone {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			errorStyle.Render(fmt.Sprintf("错误: %v", m.err)),
			"",
			statusStyle.Render("按 Enter 返回主菜单"),
		)
	}

	switch m.state {
	case removeStateSelect:
		if len(m.clones) == 0 {
			return lipgloss.JoinVertical(lipgloss.Left,
				title,
				"",
				infoStyle.Render("暂无分身。"),
				"",
				statusStyle.Render("按 Enter 返回主菜单"),
			)
		}

		items := ""
		cursor := "  "
		if m.all {
			cursor = "▸ "
		}
		items += menuItemActiveStyle.Render(cursor+"删除所有分身") + "\n"

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
			"选择要删除的分身:",
			"",
			items,
			"",
			statusStyle.Render("↑/↓ 选择  •  a 切换全部  •  Enter 确认  •  Esc 返回"),
		)

	case removeStateConfirm:
		target := "所有分身"
		if !m.all && len(m.clones) > 0 {
			target = m.clones[m.cursor].Name
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			warningStyle.Render(fmt.Sprintf("确定要删除 %s 吗？", target)),
			warningStyle.Render("此操作不可恢复!"),
			"",
			confirmStyle.Render("确认删除？(Y/N)"),
		)

	case removeStateWorking:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " 正在删除..."),
		)

	case removeStateDone:
		content := title + "\n\n"
		if m.err != nil {
			content += errorStyle.Render(fmt.Sprintf("错误: %v", m.err))
		} else {
			for _, r := range m.results {
				switch r.Status {
				case "removed":
					content += successStyle.Render(fmt.Sprintf("  ✓ 分身 %d 已删除", r.ID)) + "\n"
				case "error":
					content += errorStyle.Render(fmt.Sprintf("  ✗ 分身 %d 删除失败: %v", r.ID, r.Err)) + "\n"
				}
			}
		}
		content += "\n" + statusStyle.Render("按 Enter 返回主菜单")
		return content
	}

	return ""
}
