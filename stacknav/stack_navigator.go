package stacknav

import tea "github.com/charmbracelet/bubbletea"

type StackNavigator struct {
	pages []tea.Model
}

func (sn StackNavigator) Init() tea.Cmd {
	return sn.pages[sn.last()].Init()
}

func (sn StackNavigator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case msgPush:
		sn.pages = append(sn.pages, msg.next)
	case msgPop:
		sn.pages = sn.pages[:sn.last()]
		if len(sn.pages) == 0 {
			cmd = tea.Quit
		}
	default:
		sn.pages[sn.last()], cmd = sn.pages[sn.last()].Update(msg)
	}

	return sn, cmd
}

func (sn StackNavigator) View() string {
	return sn.pages[sn.last()].View()
}

func (sn StackNavigator) last() int {
	return len(sn.pages) - 1
}

type msgPush struct {
	next tea.Model
}

type msgPop struct{}

func New(initial tea.Model) StackNavigator {
	return StackNavigator{[]tea.Model{initial}}
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
