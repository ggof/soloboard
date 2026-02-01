package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"soloboard/db"
	"soloboard/model"
	"soloboard/page"
	"soloboard/sighandler"
	"soloboard/stacknav"
	"soloboard/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/ggof/argparse"
	"github.com/oklog/ulid/v2"

	"github.com/muesli/termenv"
)

func main() {
	p := argparse.NewParser("soloboard", "a simple TUI/CLI tool to manage your personal projects, KanBan-style.")

	debug := p.Flag("d", "debug", &argparse.Options{Default: false})

	cmdSeed := p.NewCommand("seed", "Create a test database (for debug purposes)")
	if err := p.Parse(os.Args); err != nil {
		fmt.Print(p.Usage(err))
		os.Exit(1)
	}

	env := NewEnv(*debug)

	if env.Debug {
		f, err := os.OpenFile("logs.err", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0664)
		if err != nil {
			panic(err)
		}

		log.SetOutput(f)
	}

	var err error
	switch {
	case cmdSeed.Happened():
		err = SeedDatabase(env)
	default:
		err = StartTUI(env)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

type Env struct {
	DB    page.Database
	Debug bool
}

func NewEnv(debug bool) *Env {
	dbpath := os.ExpandEnv("$HOME/.local/share/soloboard")
	if err := os.MkdirAll(dbpath, 0755); err != nil {
		panic(err)
	}

	dbfilename := path.Join(dbpath, "boards.db")

	return &Env{DB: db.NewBoardDatabase(dbfilename), Debug: debug}

}

func StartTUI(e *Env) error {
	reset := setBaseColors()
	defer reset()

	p := tea.NewProgram(
		sighandler.New(
			stacknav.New(
				page.SelectBoard(e.DB),
			),
		),
		tea.WithoutCatchPanics())
	_, err := p.Run()

	return err
}

func SeedDatabase(env *Env) error {
	boards, err := env.DB.Read()
	if err != nil {
		return err
	}


	seed := model.Board{
		ID:   ulid.Make().String(),
		Name: "seed",
		Sections: []model.Section{
			{
				ID:   ulid.Make().String(),
				Name: "TODO",
				Tasks: []model.Task{
					model.NewTask("t00", "lorem"),
					model.NewTask("t01", "lorem"),
					model.NewTask("t02", "lorem"),
					model.NewTask("t03", "lorem"),
					model.NewTask("t04", "lorem"),
					model.NewTask("t05", "lorem"),
					model.NewTask("t06", "lorem"),
					model.NewTask("t07", "lorem"),
					model.NewTask("t08", "lorem"),
					model.NewTask("t09", "lorem"),
				},
			},
			{
				ID:   ulid.Make().String(),
				Name: "IN PROGRESS",
				Tasks: []model.Task{
					model.NewTask("t10", "lorem"),
				},
			},
			{
				ID:   ulid.Make().String(),
				Name: "DONE",
				Tasks: []model.Task{
					model.NewTask("t20", "lorem"),
					model.NewTask("t21", "lorem"),
					model.NewTask("t22", "lorem"),
					model.NewTask("t23", "lorem"),
					model.NewTask("t24", "lorem"),
					model.NewTask("t25", "lorem"),
					model.NewTask("t26", "lorem"),
					model.NewTask("t27", "lorem"),
					model.NewTask("t28", "lorem"),
					model.NewTask("t29", "lorem"),
				},
			},
		},
	}

	log.Println("read database completed.")

	for i, b := range boards {
		if b.Name == seed.Name {
			boards = slices.Delete(boards, i, i+1)
			log.Println("removed old seed database.")
		}
	}

	boards = append(boards, seed)

	return env.DB.Write(boards)
}

func setBaseColors() func() {
	o := termenv.NewOutput(os.Stdout)
	bg := termenv.BackgroundColor()
	fg := termenv.ForegroundColor()

	o.SetBackgroundColor(termenv.ANSIBlack)
	o.SetCursorColor(termenv.ANSIWhite)
	o.SetForegroundColor(termenv.ANSIWhite)
	o.ClearScreen()

	return func() {
		o.SetBackgroundColor(bg)
		o.SetForegroundColor(fg)
		o.SetCursorColor(fg)
		o.ClearScreen()
	}
}
