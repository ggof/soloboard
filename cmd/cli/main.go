package main

import (
	"context"
	"math/rand/v2"
	"os"
	"soloboard/color"
	"soloboard/viewport"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "test-viewport",
				Action: TryViewport,
			},
		},
	}

	cmd.Run(context.Background(), os.Args)
}

type TryViewportApp struct {
	items []string
	viewport.Viewport
	w, h int
}

// Init implements [tea.Model].
func (t TryViewportApp) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (t TryViewportApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.w = msg.Width
		t.h = msg.Height
		t.SetSize(t.h)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return t, tea.Quit
		case "j":
			t.Prev()
		case "k":
			t.Next()
		case "n":
			n := rand.Int64N(16)
			t.items = append(t.items, strconv.FormatInt(n, 16))
			t.SetLen(len(t.items))
		}
	}

	return t, nil
}

// View implements [tea.Model].
func (t TryViewportApp) View() string {
	container := lipgloss.NewStyle().Padding(1).Margin(1).Border(lipgloss.RoundedBorder()).BorderForeground(color.Lime).Align(lipgloss.Center)
	box := container.UnsetMargins().MarginBottom(1)
	selected := box.BorderForeground(color.Orange)

	var items []string

	for i := range t.Window() {
		style := box
		if i == t.I {
			style = selected
		}

		items = append(items, style.Render(t.items[i]))
	}

	return lipgloss.Place(t.w, t.h, lipgloss.Center, lipgloss.Center, container.Render(lipgloss.JoinVertical(lipgloss.Center, items...)))
}

func TryViewport(ctx context.Context, c *cli.Command) error {
	items := make([]string, 16)
	for i := range items {
		items[i] = strconv.FormatInt(int64(i), 16)
	}

	vp := viewport.New(5)

	tea.NewProgram(TryViewportApp{Viewport: vp, items: items}).Run()

	return nil
}
