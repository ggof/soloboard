package main

import (
	"context"
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
	"github.com/oklog/ulid/v2"

	// "github.com/charmbracelet/x/term"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v3"
)

type Env struct {
	DB page.Database
}

func NewEnv() *Env {
	dbpath := os.ExpandEnv("$HOME/.local/share/soloboard")
	if err := os.MkdirAll(dbpath, 0755); err != nil {
		panic(err)
	}

	dbfilename := path.Join(dbpath, "boards.db")

	return &Env{DB: db.NewBoardDatabase(dbfilename)}

}

func StartTUI(e *Env) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
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

		// log.SetOutput(os.Stderr)

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
}

func ListBoards(e *Env) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		boards, err := e.DB.Read()
		if err != nil {
			return err
		}

		w, _, err := term.GetSize(uintptr(os.Stdout.Fd()))
		if err != nil {
			return err
		}

		size := w / 2

		parts := make([]string, len(boards)+1)
		parts[0] = lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.PlaceHorizontal(size, lipgloss.Left, "ID"),
			lipgloss.PlaceHorizontal(size, lipgloss.Left, "NAME"),
		)
		for i, b := range boards {
			lipgloss.MarkdownBorder()
			parts[i+1] = lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.PlaceHorizontal(size, lipgloss.Left, utils.EllipsisEnd(b.ID, size)),
				lipgloss.PlaceHorizontal(size, lipgloss.Left, utils.EllipsisEnd(b.Name, size)),
			)
		}

		fmt.Println(lipgloss.JoinVertical(lipgloss.Top, parts...))

		return nil
	}
}

func SeedDatabase(env *Env) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
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

		for i, b := range boards {
			if b.Name == seed.Name {
				boards = slices.Delete(boards, i, i+1)
			}
		}

		boards = append(boards, seed)

		return env.DB.Write(boards)
	}
}

func main() {
	env := NewEnv()

	program := cli.Command{
		Action: StartTUI(env),
		Commands: []*cli.Command{
			{
				Name:        "list",
				Description: "List the current boards",
				Action:      ListBoards(env),
			},
			{
				Name:        "seed",
				Description: "Create a test database (for debug purposes)",
				Action:      SeedDatabase(env),
			},
		},
	}

	f, err := os.OpenFile("logs.err", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	log.SetOutput(f)

	if err := program.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
