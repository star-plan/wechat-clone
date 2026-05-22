package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deali/wechat-clone/internal/app"
)

// Page represents the current TUI page.
type Page int

const (
	PageMenu Page = iota
	PageCreate
	PageList
	PageOpen
	PageRepair
	PageRemove
	PageDoctor
	PageQuit
)

// menuItems defines the main menu options.
var menuItems = []struct {
	label       string
	description string
	page        Page
}{
	{"创建分身", "创建指定数量的微信分身", PageCreate},
	{"查看分身列表", "查看当前已有的所有分身", PageList},
	{"启动分身", "打开指定或所有分身", PageOpen},
	{"修复分身", "重新修复分身签名和 quarantine", PageRepair},
	{"删除分身", "删除指定或所有分身", PageRemove},
	{"环境检查", "检查系统环境是否满足要求", PageDoctor},
}

// MainModel is the top-level TUI model that routes between pages.
type MainModel struct {
	app       *app.App
	page      Page
	cursor    int
	width     int
	height    int
	subModel  tea.Model
	spinner   spinner.Model
}

// NewMainModel creates a new MainModel.
func NewMainModel(a *app.App) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primary)

	return MainModel{
		app:     a,
		page:    PageMenu,
		cursor:  0,
		spinner: s,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.page == PageMenu {
				return m, tea.Quit
			}
			// Go back to menu from sub-pages
			m.page = PageMenu
			m.subModel = nil
			return m, nil
		case "esc":
			if m.page != PageMenu {
				m.page = PageMenu
				m.subModel = nil
				return m, nil
			}
		}
	}

	// Route to sub-page if active
	if m.subModel != nil {
		var cmd tea.Cmd
		m.subModel, cmd = m.subModel.Update(msg)

		// Check if sub-model wants to return to menu
		if sm, ok := m.subModel.(SubModel); ok && sm.Done() {
			m.page = PageMenu
			m.subModel = nil
			return m, nil
		}

		return m, cmd
	}

	// Menu page handling
	if m.page == PageMenu {
		return m.updateMenu(msg)
	}

	return m, nil
}

func (m MainModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}
		case "enter":
			item := menuItems[m.cursor]
			m.page = item.page
			m.subModel = m.createSubModel(item.page)
			return m, m.subModel.Init()
		}
	}
	return m, nil
}

func (m MainModel) createSubModel(page Page) tea.Model {
	switch page {
	case PageCreate:
		return NewCreateModel(m.app)
	case PageList:
		return NewListModel(m.app)
	case PageOpen:
		return NewOpenModel(m.app)
	case PageRepair:
		return NewRepairModel(m.app)
	case PageRemove:
		return NewRemoveModel(m.app)
	case PageDoctor:
		return NewDoctorModel(m.app)
	default:
		return nil
	}
}

func (m MainModel) View() string {
	if m.page == PageMenu {
		return m.viewMenu()
	}

	if m.subModel != nil {
		return m.subModel.View()
	}

	return ""
}

func (m MainModel) viewMenu() string {
	title := titleStyle.Render("微信分身管理工具")
	subtitle := infoStyle.Render("使用 ↑↓ 键选择，Enter 确认，q 退出")

	var items string
	for i, item := range menuItems {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "▸ "
			style = menuItemActiveStyle
		}
		items += style.Render(cursor+item.label) + "\n"
	}

	help := statusStyle.Render("↑/↓ 移动  •  Enter 选择  •  q 退出")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		subtitle,
		"",
		items,
		help,
	)
}

// SubModel is an interface for sub-page models that can signal completion.
type SubModel interface {
	tea.Model
	Done() bool
}
