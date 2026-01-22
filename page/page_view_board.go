package page

import (
	"slices"
	"soloboard/color"
	"soloboard/model"
	"soloboard/viewport"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sectionWidth = 40

type PageViewBoard struct {
	viewport.Viewport
	boards  []model.Board
	columns []viewport.Viewport
	db      Database

	i, w, h int
}

func ViewBoard(db Database, boards []model.Board, i, width, height int) PageViewBoard {
	sw := max(40, width/3)

	p := PageViewBoard{
		Viewport: viewport.New(sw),
		boards:   boards,
		db:       db,
		i:        i,
		w:        width,
		h:        height,
	}

	p.SetSize(p.w)
	p.SetLen(len(p.boards[i].Sections))

	p.columns = make([]viewport.Viewport, len(p.boards[i].Sections))
	for j := range boards[i].Sections {
		p.columns[j] = viewport.New(2 + 3)    // allow 3 lines of text per task
		p.columns[j].SetSize(p.h - 3 - 3 - 2) // 3 lines for title, 3 lines for col title, 2 lines for borders
		p.columns[j].SetLen(len(boards[i].Sections[j].Tasks))
	}

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
		for j := range p.columns {
			p.columns[j].SetSize(p.h - 3 - 3 - 2)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			p.Prev()
		case "l":
			p.Next()
		case "j":
			p.columns[p.I].Next()
		case "k":
			p.columns[p.I].Prev()
		case "H":
			oldcol := p.I
			p.moveTask(p.I, p.I-1, p.columns[p.I].I)
			p.columns[p.I].SetLen(len(p.boards[p.i].Sections[p.I].Tasks))
			if p.columns[p.I].I == len(p.boards[p.i].Sections[p.I].Tasks) {
				p.columns[p.I].I--
			}
			p.Prev()
			if oldcol != p.I {
				p.columns[p.I].I = len(p.boards[p.i].Sections[p.I].Tasks) - 1
				p.columns[p.I].SetLen(len(p.boards[p.i].Sections[p.I].Tasks))
			}
			cmd = p.saveBoards()
		case "L":
			oldcol := p.I
			p.moveTask(p.I, p.I+1, p.columns[p.I].I)
			if p.columns[p.I].I == len(p.boards[p.i].Sections[p.I].Tasks) {
				p.columns[p.I].I--
			}
			p.columns[p.I].SetLen(len(p.boards[p.i].Sections[p.I].Tasks))
			p.Next()
			if oldcol != p.I {
				p.columns[p.I].I = len(p.boards[p.i].Sections[p.I].Tasks) - 1
				p.columns[p.I].SetLen(len(p.boards[p.i].Sections[p.I].Tasks))
			}

			cmd = p.saveBoards()
		case "K":
			p.swapTask(p.I, p.columns[p.I].I, p.columns[p.I].I-1)
			p.columns[p.I].Prev()
			cmd = p.saveBoards()
		case "J":
			p.swapTask(p.I, p.columns[p.I].I, p.columns[p.I].I+1)
			p.columns[p.I].Next()
			cmd = p.saveBoards()
		}
	}
	return p, cmd
}

func (p PageViewBoard) View() string {
	d := lipgloss.NewStyle().Width(sectionWidth).Height(p.h - 3 - 2).AlignVertical(lipgloss.Top).Border(lipgloss.RoundedBorder())

	s := d.BorderForeground(color.Lime)

	taskBox := lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.RoundedBorder())
	selectedTask := taskBox.BorderForeground(color.LightBlue)

	var cols []string
	for i := range p.Window() {
		vp := p.columns[i]
		var col []string

		colTitleStyle := lipgloss.NewStyle().Padding(1).Underline(true)
		colTitle := lipgloss.PlaceHorizontal(sectionWidth-2, lipgloss.Center, colTitleStyle.Render(p.boards[p.i].Sections[i].Name))

		col = append(col, colTitle)
		for j := range vp.Window() {
			ss := taskBox
			if p.I == i && vp.I == j {
				ss = selectedTask
			}

			task := lipgloss.Place(sectionWidth-2, 3, lipgloss.Left, lipgloss.Center, p.boards[p.i].Sections[i].Tasks[j].Name)
			task = ss.Render(task)

			col = append(col, task)
		}
		ss := d
		if p.I == i {
			ss = s
		}

		cols = append(cols, ss.Render(lipgloss.JoinVertical(lipgloss.Center, col...)))
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.Place(p.w, 3, lipgloss.Center, lipgloss.Center, lipgloss.NewStyle().Bold(true).Underline(true).Render(p.boards[p.i].Name)),
		lipgloss.PlaceHorizontal(p.w, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, cols...)),
	)
}

func (p *PageViewBoard) swapTask(s, i, j int) {
	if j < 0 || j == len(p.boards[p.i].Sections[s].Tasks) {
		return
	}

	p.boards[p.i].Sections[s].Tasks[i], p.boards[p.i].Sections[s].Tasks[j] = p.boards[p.i].Sections[s].Tasks[j], p.boards[p.i].Sections[s].Tasks[i]
}

func (p *PageViewBoard) moveTask(si, sj, i int) {
	if sj < 0 || sj == len(p.boards[p.i].Sections) {
		return
	}

	task := p.boards[p.i].Sections[si].Tasks[i]

	p.boards[p.i].Sections[si].Tasks = slices.Delete(p.boards[p.i].Sections[si].Tasks, i, i+1)
	p.boards[p.i].Sections[sj].Tasks = append(p.boards[p.i].Sections[sj].Tasks, task)
}

func (p PageViewBoard) saveBoards() tea.Cmd {
	return func() tea.Msg {
		if err := p.db.Write(p.boards); err != nil {
			panic(err)
		}

		return nil
	}
}
