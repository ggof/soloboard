package page

import (
	"soloboard/color"
	"soloboard/model"
	"soloboard/viewport"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sectionWidth = 40

type PageViewBoard struct {
	viewport.Viewport
	board model.Board

	w, h int
}

func ViewBoard(board model.Board, width, height int) PageViewBoard {
	p := PageViewBoard{
		Viewport: viewport.New(sectionWidth),
		board:    board,
		w:        width,
		h:        height,
	}

	p.SetSize(p.w)
	p.SetLen(len(p.board.Sections))

	return p
}

func (p PageViewBoard) Init() tea.Cmd {
	return nil
}

func (p PageViewBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.w, p.h = msg.Width, msg.Height
		p.SetSize(p.w)
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			p.Prev()
		case "l":
			p.Next()
		}
	}
	return p, cmd
}

func (p PageViewBoard) View() string {
	title := lipgloss.PlaceVertical(3, lipgloss.Center, p.board.Name)
	d := lipgloss.NewStyle().Width(sectionWidth).Height(p.h-5).Padding(0, 1).Align(lipgloss.Center, lipgloss.Top).Border(lipgloss.RoundedBorder())

	s := d.BorderForeground(color.Lime)

	var cols []string
	for i := range p.Window() {
		if p.I == i {
			cols = append(cols, s.Render(p.board.Sections[i].Name))
		} else {
			cols = append(cols, d.Render(p.board.Sections[i].Name))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.PlaceHorizontal(p.w, lipgloss.Center, title),
		lipgloss.PlaceHorizontal(p.w, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, cols...)),
	)
}
