package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primary   = lipgloss.Color("#07C160") // WeChat green
	secondary = lipgloss.Color("#576B95")
	dim       = lipgloss.Color("#999999")
	danger    = lipgloss.Color("#FA5151")
	warning   = lipgloss.Color("#FFC300")
	success   = lipgloss.Color("#07C160")

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primary).
			MarginBottom(1)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)

	menuItemActiveStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(2).
				Foreground(primary).
				Bold(true)

	// Info text
	infoStyle = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true)

	// Success message
	successStyle = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	// Error message
	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)

	// Warning message
	warningStyle = lipgloss.NewStyle().
			Foreground(warning)

	// Table header
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primary).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder())

	// Table cell
	tableCellStyle = lipgloss.NewStyle().
			PaddingRight(2)

	// Status bar at bottom
	statusStyle = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true).
			MarginTop(1)

	// Box container
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondary).
			Padding(1, 2)

	// Confirm dialog
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warning).
			Padding(1, 2).
			MarginTop(1)
)
