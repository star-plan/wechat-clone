package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
	"github.com/deali/wechat-clone/internal/macos"
)

type ListModel struct {
	app    *app.App
	clones []macos.CloneInfo
	err    error
	done   bool
}

func NewListModel(a *app.App) *ListModel {
	return &ListModel{app: a}
}

func (m *ListModel) Init() tea.Cmd {
	return func() tea.Msg {
		clones, err := m.app.ListClones()
		return listLoadedMsg{clones: clones, err: err}
	}
}

func (m *ListModel) Done() bool {
	return m.done
}

type listLoadedMsg struct {
	clones []macos.CloneInfo
	err    error
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			m.done = true
			return m, nil
		}
	case listLoadedMsg:
		m.clones = msg.clones
		m.err = msg.err
	}
	return m, nil
}

func (m *ListModel) View() string {
	title := titleStyle.Render("微信分身列表")

	if m.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			errorStyle.Render(fmt.Sprintf("获取列表失败: %v", m.err)),
			"",
			statusStyle.Render("按 Enter 返回主菜单"),
		)
	}

	if len(m.clones) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			infoStyle.Render("暂无分身。请先创建分身。"),
			"",
			statusStyle.Render("按 Enter 返回主菜单"),
		)
	}

	// Table header
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		tableHeaderStyle.Width(8).Render("编号"),
		tableHeaderStyle.Width(28).Render("名称"),
		tableHeaderStyle.Width(36).Render("Bundle ID"),
	)

	// Table rows
	var rows string
	for _, c := range m.clones {
		row := lipgloss.JoinHorizontal(lipgloss.Left,
			tableCellStyle.Width(8).Render(fmt.Sprintf("%d", c.ID)),
			tableCellStyle.Width(28).Render(c.Name),
			tableCellStyle.Width(36).Render(c.BundleID),
		)
		rows += row + "\n"
	}

	count := infoStyle.Render(fmt.Sprintf("共 %d 个分身", len(m.clones)))

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		header,
		rows,
		count,
		"",
		statusStyle.Render("按 Enter 返回主菜单"),
	)
}
