package main

import (
	"os"
	"path"
	"soloboard/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type BoardDatabase interface {
	Read() ([]Board, error)
	Write([]Board) error
}

type Task struct {
	ID          int
	Name        string
	Description string
}

type Section struct {
	ID    int
	Name  string
	Tasks []Task
}

type Board struct {
	ID       int
	Name     string
	Sections []Section
}

func main() {
	dbpath := os.ExpandEnv("$HOME/.local/share/soloboard")
	if err := os.MkdirAll(dbpath, 0755); err != nil {
		panic(err)
	}

	dbfilename := path.Join(dbpath, "boards.db")

	o := termenv.NewOutput(os.Stdout)

	bg := termenv.BackgroundColor()
	fg := termenv.ForegroundColor()

	defer func() {
		o.SetBackgroundColor(bg)
		o.SetForegroundColor(fg)
		o.SetCursorColor(fg)
		o.ClearScreen()
	}()

	o.SetBackgroundColor(termenv.ANSIBlack)
	o.SetCursorColor(termenv.ANSIWhite)
	o.SetForegroundColor(termenv.ANSIWhite)

	p := tea.NewProgram(PageSelectBoard{db: NewBoardDatabase(dbfilename), vp: components.NewViewport(0, 5)}, tea.WithoutCatchPanics())
	_, err := p.Run()
	if err != nil {
		panic(err)
	}

}
