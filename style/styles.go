package style

import (
	"soloboard/color"

	"github.com/charmbracelet/lipgloss"
)

var Box = lipgloss.NewStyle().
	Align(lipgloss.Center, lipgloss.Center).
	Border(lipgloss.RoundedBorder()).
	Padding(1)

var SelectedBox = Box.BorderForeground(color.LightBlue)

var Column = lipgloss.NewStyle().
	AlignVertical(lipgloss.Top).
	Border(lipgloss.RoundedBorder())

var SelectedColumn = Column.BorderForeground(color.Lime)
