package sighandler

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	inner tea.Model
}

func New(inner tea.Model) Model {
	return Model{inner}

}

func (m Model) Init() tea.Cmd {
	return m.inner.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.inner.View()
}
