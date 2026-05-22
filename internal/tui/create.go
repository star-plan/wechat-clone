package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
)

type createState int

const (
	createStateInput createState = iota
	createStateConfirm
	createStateWorking
	createStateDone
)

type createDoneMsg struct {
	results []app.CreateResult
	err     error
}

type CreateModel struct {
	app       *app.App
	state     createState
	textInput textinput.Model
	spinner   spinner.Model
	force     bool
	count     int
	results   []app.CreateResult
	err       error
	done      bool
}

func NewCreateModel(a *app.App) *CreateModel {
	ti := textinput.New()
	ti.Placeholder = "输入分身数量 (例如: 3)"
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return &CreateModel{
		app:       a,
		state:     createStateInput,
		textInput: ti,
		spinner:   s,
	}
}

func (m *CreateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *CreateModel) Done() bool {
	return m.done
}

func (m *CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case createStateInput:
			switch msg.String() {
			case "enter":
				val := m.textInput.Value()
				count, err := strconv.Atoi(val)
				if err != nil || count <= 0 {
					m.err = fmt.Errorf("请输入有效的正整数")
					return m, nil
				}
				m.count = count
				m.err = nil
				m.state = createStateConfirm
				return m, nil
			}
		case createStateConfirm:
			switch msg.String() {
			case "y", "Y":
				m.state = createStateWorking
				return m, tea.Batch(m.spinner.Tick, m.doCreate())
			case "n", "N", "esc":
				m.done = true
				return m, nil
			}
		case createStateDone:
			switch msg.String() {
			case "enter":
				m.done = true
				return m, nil
			}
		}

	case createDoneMsg:
		m.results = msg.results
		m.err = msg.err
		m.state = createStateDone
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *CreateModel) doCreate() tea.Cmd {
	return func() tea.Msg {
		results, err := m.app.CreateClones(m.count, m.force)
		return createDoneMsg{results: results, err: err}
	}
}

func (m *CreateModel) View() string {
	title := titleStyle.Render("创建微信分身")

	switch m.state {
	case createStateInput:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			"请输入要创建的分身数量:",
			"",
			m.textInput.View(),
			"",
			infoStyle.Render("Enter 确认  •  Esc 返回"),
		)

	case createStateConfirm:
		msg := fmt.Sprintf("即将创建 %d 个微信分身", m.count)
		if m.err != nil {
			msg = errorStyle.Render(m.err.Error())
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			msg,
			"",
			confirmStyle.Render("确认创建？(Y/N)"),
		)

	case createStateWorking:
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.JoinHorizontal(lipgloss.Center, m.spinner.View(), " 正在创建分身..."),
		)

	case createStateDone:
		content := title + "\n\n"
		if m.err != nil {
			content += errorStyle.Render(fmt.Sprintf("错误: %v", m.err))
		} else {
			for _, r := range m.results {
				switch r.Status {
				case "created":
					content += successStyle.Render(fmt.Sprintf("  ✓ 分身 %d 创建成功", r.ID)) + "\n"
				case "skipped":
					content += warningStyle.Render(fmt.Sprintf("  - 分身 %d 已存在，跳过", r.ID)) + "\n"
				case "error":
					content += errorStyle.Render(fmt.Sprintf("  ✗ 分身 %d 创建失败: %v", r.ID, r.Err)) + "\n"
				}
			}
		}
		content += "\n" + statusStyle.Render("按 Enter 返回主菜单")
		return content
	}

	return ""
}
