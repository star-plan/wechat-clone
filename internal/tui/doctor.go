package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
)

type doctorState int

const (
	doctorStateWorking doctorState = iota
	doctorStateDone
)

type doctorDoneMsg struct {
	results []app.CheckResult
}

type DoctorModel struct {
	app     *app.App
	state   doctorState
	spinner spinner.Model
	results []app.CheckResult
	done    bool
}

func NewDoctorModel(a *app.App) *DoctorModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return &DoctorModel{
		app:     a,
		state:   doctorStateWorking,
		spinner: s,
	}
}

func (m *DoctorModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.doDoctor())
}

func (m *DoctorModel) Done() bool {
	return m.done
}

func (m *DoctorModel) doDoctor() tea.Cmd {
	return func() tea.Msg {
		results := m.app.Doctor()
		return doctorDoneMsg{results: results}
	}
}

func (m *DoctorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case doctorStateDone:
			switch msg.String() {
			case "enter":
				m.done = true
				return m, nil
			}
		}

	case doctorDoneMsg:
		m.results = msg.results
		m.state = doctorStateDone
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *DoctorModel) View() string {
	title := titleStyle.Render("环境检查")

	switch m.state {
	case doctorStateWorking:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " 正在检查环境..."),
		)

	case doctorStateDone:
		content := title + "\n\n"
		allOK := true
		for _, r := range m.results {
			var icon string
			var style lipgloss.Style
			switch r.Status {
			case app.CheckOK:
				icon = "✓"
				style = successStyle
			case app.CheckWarning:
				icon = "⚠"
				style = warningStyle
				allOK = false
			case app.CheckError:
				icon = "✗"
				style = errorStyle
				allOK = false
			}
			content += fmt.Sprintf("  %s %s: %s\n", style.Render(icon), r.Name, r.Message)
		}

		content += "\n"
		if allOK {
			content += successStyle.Render("环境检查通过! 可以正常使用。")
		} else {
			content += warningStyle.Render("存在以上问题，请修复后再使用。")
		}

		content += "\n\n" + statusStyle.Render("按 Enter 返回主菜单")
		return content
	}

	return ""
}
