package main

import (
	"log"
	"os"
	"path"
	"soloboard/db"
	"soloboard/page"
	"soloboard/sighandler"
	"soloboard/stacknav"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

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

	log.SetOutput(os.Stderr)

	p := tea.NewProgram(sighandler.New(stacknav.New(page.SelectBoard(db.NewBoardDatabase(dbfilename)))), tea.WithoutCatchPanics())
	_, err := p.Run()
	if err != nil {
		panic(err)
	}

}
