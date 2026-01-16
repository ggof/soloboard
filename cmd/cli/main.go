package main

import (
	"context"
	"math/rand/v2"
	"os"
	"soloboard/components"
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
	items []Item
	vp    *components.Viewport
}

// Init implements [tea.Model].
func (t TryViewportApp) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (t TryViewportApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return t, tea.Quit
		case "j":
			t.vp.Down()
		case "k":
			t.vp.Up()
		case "n":
			n := rand.Int64N(16)
			t.items = append(t.items, Item{strconv.FormatInt(n, 16)})
			t.vp.SetList(RenderAsBox(t.items))

		}

	case tea.WindowSizeMsg:
		t.vp.SetSize(msg.Width, msg.Height)
	}

	return t, nil
}

// View implements [tea.Model].
func (t TryViewportApp) View() string {
	return t.vp.Render()
}

var tva tea.Model = (*TryViewportApp)(nil)

type Item struct {
	Title string
}

func TryViewport(ctx context.Context, c *cli.Command) error {
	items := make([]Item, 16)
	for i := range items {
		items[i].Title = strconv.FormatInt(int64(i), 16)
	}

	vp := components.NewViewport(RenderAsBox(items), len(items))

	tea.NewProgram(TryViewportApp{items, vp}).Run()

	return nil
}

type itemSliceRenderer struct {
	items         []Item
	defaultStyle  lipgloss.Style
	selectedStyle lipgloss.Style
}

func (isr itemSliceRenderer) RenderItem(i int, selected bool) string {
	s := isr.defaultStyle
	if selected {
		s = isr.selectedStyle
	}
	return s.Render(isr.items[i].Title)
}

func (isr itemSliceRenderer) Len() int { return len(isr.items) }

func RenderAsBox(items []Item) itemSliceRenderer {
	defaultStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).MarginBottom(1).Padding(1, 0).Width(40).Align(lipgloss.Center)
	selectedStyle := defaultStyle.BorderForeground(lipgloss.Color("#ff0000"))

	return itemSliceRenderer{
		items, defaultStyle, selectedStyle,
	}
}
