package page

import (
	"slices"
	"soloboard/color"
	"soloboard/model"
	"soloboard/stacknav"
	"soloboard/utils"
	"soloboard/viewport"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Database interface {
	Read() ([]model.Board, error)
	Write([]model.Board) error
}

func SelectBoard(db Database) PageSelectBoard {
	return PageSelectBoard{db: db, Viewport: viewport.New(5)}
}

type PageSelectBoard struct {
	boards []model.Board
	viewport.Viewport

	title      string
	insertMode bool
	db         Database

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
	case []model.Board:
		p.boards = msg
		p.SetLen(len(p.boards) + 1)
	case tea.KeyMsg:
		key := msg.String()
		switch key {
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
		p.SetSize(2 * p.h / 3)
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
		if p.I == len(p.boards) {
			p.boards = append(p.boards, model.NewBoard(p.title))
			p.SetLen(len(p.boards) + 1)
		} else {
			p.boards[p.I].Name = p.title
		}
		p.title = ""
		p.insertMode = false
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
		p.Next()
	case "k":
		p.Prev()
	case "r":
		if p.I == len(p.boards) {
			// `new board` needs `enter` to edit
			break
		}
		p.insertMode = true
		p.title = p.boards[p.I].Name
	case "d", "x":
		if p.I < len(p.boards) {
			p.boards = slices.Delete(p.boards, p.I, p.I+1)
			p.SetLen(len(p.boards) + 1)
			cmd = p.saveBoards()
		}
	case "enter":
		if p.I == len(p.boards) {
			p.insertMode = true
		} else {
			cmd = stacknav.Push(ViewBoard(p.boards, p.I, p.w, p.h))
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

	s := d.BorderForeground(color.LightBlue)

	var elems []string

	for i := range p.Window() {
		ss := d
		if p.I == i {
			ss = s
			if p.insertMode {
				ss = ss.BorderForeground(color.Purple)
			}
		}

		var text string

		switch {
		case p.insertMode && p.I == i:
			ss = ss.Align(lipgloss.Left)
			text = utils.EllipsisBeg(p.title, w-3) + "|"
		case i == len(p.boards):
			text = utils.EllipsisEnd("New Board", w-2)
		default:
			text = utils.EllipsisEnd(p.boards[i].Name, w-2)
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
