package main

import (
	"slices"
	"soloboard/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PageSelectBoard struct {
	boards []Board
	vp     *components.Viewport

	title      string
	insertMode bool
	db         BoardDatabase

	w int
	h int
}

func (p PageSelectBoard) Init() tea.Cmd {
	return func() tea.Msg {
		// TODO: fetch boards from sqlite
		boards, err := p.db.Read()
		if err != nil {
			panic(err)
		}
		return boards
	}
}

func (p PageSelectBoard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []Board:
		p.boards = msg
		p.vp.SetList(Boxed(p.boards))
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c":
			return p, tea.Quit
		default:
			if !p.insertMode {
				return p.handleNormalMode(key)
			} else {
				return p.handleInsertMode(key)
			}
		}
	case tea.WindowSizeMsg:
		p.w = msg.Width
		p.h = msg.Height
		p.vp.SetSize(p.w, p.h)
	}
	return p, nil
}

func (p PageSelectBoard) handleInsertMode(key string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key == "esc":
		p.title = ""
		p.insertMode = false
	case len(key) == 1:
		p.title += key

	case key == "enter":
		p.boards = append(p.boards, Board{Name: p.title, Sections: []Section{{Name: "TODO"}, {Name: "IN PROGRESS"}, {Name: "DONE"}}})
		p.title = ""
		p.insertMode = false
		p.vp.SetList(Boxed(p.boards))
		cmd = p.saveBoards()
	case key == "backspace":
		if p.title != "" {
			p.title = p.title[:len(p.title)-1]
		}
	}

	return p, cmd
}

func (p PageSelectBoard) handleNormalMode(key string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch key {
	case "j":
		p.vp.Down()
	case "k":
		p.vp.Up()
	case "d", "x":
		if p.vp.I < len(p.boards) {
			p.boards = slices.Delete(p.boards, p.vp.I, p.vp.I+1)
			p.vp.SetList(Boxed(p.boards))
			cmd = p.saveBoards()
		}
	case "enter":
		if p.vp.I == len(p.boards) {
			p.insertMode = true
		} else {
			// TODO: return the next page
		}
	}

	return p, cmd
}

func (p PageSelectBoard) View() string {
	return lipgloss.Place(p.w, p.h, lipgloss.Center, lipgloss.Center, p.vp.Render())
}

func (p PageSelectBoard) saveBoards() tea.Cmd {
	return func() tea.Msg {
		if err := p.db.Write(p.boards); err != nil {
			panic(err)
		}

		return nil
	}
}

type boardRenderer struct {
	boards        []Board
	defaultStyle  lipgloss.Style
	selectedStyle lipgloss.Style
}

func (b boardRenderer) Len() int { return len(b.boards) + 1 } // account for the "new board" box

func (b boardRenderer) RenderItem(i int, selected bool) string {
	s := b.defaultStyle
	if selected {
		s = b.selectedStyle
	}
	if i == len(b.boards) {
		return s.Render("New board")
	}

	return s.Render(b.boards[i].Name)
}

func Boxed(boards []Board) boardRenderer {
	d := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		MarginTop(1).
		Width(40)

	s := d.BorderForeground(ColorLightBlue)
	return boardRenderer{boards, d, s}
}
