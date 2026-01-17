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
		p.vp.SetLen(len(p.boards) + 1)
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
		p.vp.SetSize(2 * p.h / 3)
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
		p.vp.SetLen(len(p.boards) + 1)
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
			p.vp.SetLen(len(p.boards) + 1)
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
	w := max(40, p.w/3)
	d := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		MarginBottom(1).
		Width(w)

	s := d.BorderForeground(ColorLightBlue)

	var elems []string

	for i := range p.vp.Window() {
		ss := d
		if p.vp.I == i {
			ss = s
			if i == len(p.boards) {
				ss = ss.BorderForeground(ColorPurple)
			}
		}

		var text string

		if i == len(p.boards) {
			if p.insertMode {
				ss = ss.Align(lipgloss.Left)
				text = ellipsisBeg(p.title, w-3) + "|"
			} else {
				text = "New board"
			}
		} else {
			text = ellipsisEnd(p.boards[i].Name, w-2)
		}

		elems = append(elems, ss.Render(text))
	}

	return lipgloss.Place(p.w, p.h, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, elems...))
}

func (p PageSelectBoard) saveBoards() tea.Cmd {
	return func() tea.Msg {
		if err := p.db.Write(p.boards); err != nil {
			panic(err)
		}

		return nil
	}
}

func ellipsisBeg(text string, w int) string {
	if len(text) > w {
		tbs := []byte(text[len(text)-w:])
		copy(tbs[0:3], "...")
		text = string(tbs)
	}
	return text
}

func ellipsisEnd(text string, w int) string {
	if len(text) > w {
		tbs := []byte(text[:w])
		copy(tbs[w-3:], "...")
		text = string(tbs)
	}
	return text
}
