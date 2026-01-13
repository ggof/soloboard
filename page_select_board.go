package main

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	// TODO: check the current state (editing or nah), update model based on key pressed
	switch msg := msg.(type) {
	case []Board:
		p.boards = msg
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c":
			return p, tea.Quit
		default:
			if !p.editing {
				return p.handleNonEditingKey(key)
			} else {
				return p.handleEdit(key)
			}
		}
	case tea.WindowSizeMsg:
		p.w = msg.Width
		p.h = msg.Height
	}
	return p, nil
}

func (p PageSelectBoard) handleEdit(key string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key == "esc":
		p.title = ""
		p.editing = false
	case len(key) == 1:
		p.title += key

	case key == "enter":
		p.boards = append(p.boards, Board{Name: p.title, Sections: []Section{{Name: "TODO"}, {Name: "IN PROGRESS"}, {Name: "DONE"}}})
		p.title = ""
		p.editing = false
		cmd = p.saveBoards()
	}

	return p, cmd
}

func (p PageSelectBoard) handleNonEditingKey(key string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch key {
	case "j":
		if p.selected < len(p.boards) {
			p.selected++
		}
	case "k":
		if p.selected > 0 {
			p.selected--
		}
	case "d", "x":
		if p.selected < len(p.boards) {
			p.boards = slices.Delete(p.boards, p.selected, p.selected+1)
			cmd = p.saveBoards()
		}
	case "enter":
		if p.selected == len(p.boards) {
			p.editing = true
		} else {
			// TODO: return the next page
		}
	}

	return p, cmd
}

func (p PageSelectBoard) View() string {
	elementStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Color0).
		Padding(1, 2).
		MarginTop(1).
		Width(40)

	names := make([]string, len(p.boards))
	for i, b := range p.boards {
		if i == p.selected {
			names[i] = elementStyle.BorderForeground(Color8).Render(b.Name)
		} else {
			names[i] = elementStyle.Render(b.Name)
		}
	}

	var lastElement string
	if p.editing {
		lastElement = elementStyle.AlignHorizontal(lipgloss.Left).BorderForeground(Color3).Render(p.title + "|")
	} else {
		if p.selected == len(p.boards) {
			lastElement = elementStyle.BorderForeground(Color3).Render("new board")
		} else {
			lastElement = elementStyle.Render("new board")
		}
	}
	names = append(names, lastElement)

	return lipgloss.Place(p.w, p.h, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, names...))
}

func (p PageSelectBoard) saveBoards() tea.Cmd {
	return func() tea.Msg {
		if err := p.db.Write(p.boards); err != nil {
			panic(err)
		}

		return nil
	}
}
