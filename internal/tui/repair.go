package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/macos"
)

type repairState int

const (
	repairStateSelect repairState = iota
	repairStateWorking
	repairStateDone
)

type repairLoadedMsg struct {
	clones []macos.CloneInfo
	err    error
}

type repairDoneMsg struct {
	results []app.RepairResult
	err     error
}

type RepairModel struct {
	app     *app.App
	state   repairState
	clones  []macos.CloneInfo
	cursor  int
	all     bool
	spinner spinner.Model
	results []app.RepairResult
	err     error
	done    bool
}

func NewRepairModel(a *app.App) *RepairModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return &RepairModel{
		app:     a,
		state:   repairStateSelect,
		all:     true,
		spinner: s,
	}
}

func (m *RepairModel) Init() tea.Cmd {
	return func() tea.Msg {
		clones, err := m.app.ListClones()
		return repairLoadedMsg{clones: clones, err: err}
	}
}

func (m *RepairModel) Done() bool {
	return m.done
}

func (m *RepairModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case repairStateSelect:
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
				m.state = repairStateWorking
				return m, tea.Batch(m.spinner.Tick, m.doRepair())
			}
		case repairStateDone:
			switch msg.String() {
			case "enter":
				m.done = true
				return m, nil
			}
		}

	case repairLoadedMsg:
		m.clones = msg.clones
		m.err = msg.err

	case repairDoneMsg:
		m.results = msg.results
		m.err = msg.err
		m.state = repairStateDone
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *RepairModel) doRepair() tea.Cmd {
	return func() tea.Msg {
		if m.all {
			results, err := m.app.RepairClones()
			return repairDoneMsg{results: results, err: err}
		}
		if len(m.clones) == 0 {
			return repairDoneMsg{err: fmt.Errorf("没有找到任何分身")}
		}
		id := m.clones[m.cursor].ID
		results, err := m.app.RepairClones(id)
		return repairDoneMsg{results: results, err: err}
	}
}

func (m *RepairModel) View() string {
	title := titleStyle.Render("修复微信分身")

	if m.err != nil && m.state != repairStateDone {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			errorStyle.Render(fmt.Sprintf("错误: %v", m.err)),
			"",
			statusStyle.Render("按 Enter 返回主菜单"),
		)
	}

	switch m.state {
	case repairStateSelect:
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
		items += menuItemActiveStyle.Render(cursor+"修复所有分身") + "\n"

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
			"选择要修复的分身:",
			"",
			items,
			"",
			statusStyle.Render("↑/↓ 选择  •  a 切换全部  •  Enter 确认  •  Esc 返回"),
		)

	case repairStateWorking:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " 正在修复..."),
		)

	case repairStateDone:
		content := title + "\n\n"
		if m.err != nil {
			content += errorStyle.Render(fmt.Sprintf("错误: %v", m.err))
		} else {
			for _, r := range m.results {
				switch r.Status {
				case "repaired":
					content += successStyle.Render(fmt.Sprintf("  ✓ 分身 %d 修复成功", r.ID)) + "\n"
				case "error":
					content += errorStyle.Render(fmt.Sprintf("  ✗ 分身 %d 修复失败: %v", r.ID, r.Err)) + "\n"
				}
			}
		}
		content += "\n" + statusStyle.Render("按 Enter 返回主菜单")
		return content
	}

	return ""
}
