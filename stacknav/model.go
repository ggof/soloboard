package stacknav

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	pages []tea.Model
}

func (m Model) Init() tea.Cmd {
	return m.pages[m.last()].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case msgPush:
		m.pages = append(m.pages, msg.next)
	case msgPop:
		m.pages = m.pages[:m.last()]
		if len(m.pages) == 0 {
			cmd = tea.Quit
		}
	default:
		m.pages[m.last()], cmd = m.pages[m.last()].Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.pages[m.last()].View()
}

func (m Model) last() int {
	return len(m.pages) - 1
}

type msgPush struct {
	next tea.Model
}

type msgPop struct{}

func New(initial tea.Model) Model {
	return Model{[]tea.Model{initial}}
}

func Push(next tea.Model) tea.Cmd {
	return func() tea.Msg {
		return msgPush{next}
	}
}

func Pop() tea.Cmd {
	return func() tea.Msg {
		return msgPop{}
	}
}
