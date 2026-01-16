package main

import (
	"github.com/charmbracelet/lipgloss"
)

type VerticalViewport struct {
	items []string
	i     int
	wb    int
	size  int
	gap   int
}

func (v *VerticalViewport) Next() {
	if v.i < len(v.items)-1 {
		v.i++
	}

	if v.i == v.wb+v.size {
		v.wb++
	}
}

func (v *VerticalViewport) Prev() {
	if v.i > 0 {
		v.i--
	}

	if v.i < v.wb {
		v.wb--
	}
}

func (v VerticalViewport) Render() string {
	if len(v.items) == 0 {
		return ""
	}

	rowStyle := lipgloss.NewStyle().MarginTop(v.gap)
	rs := []string{v.items[v.wb]}
	for i := v.wb + 1; i < len(v.items) && i < v.wb+v.size; i++ {
		rs = append(rs, rowStyle.Render(v.items[i]))
	}

	return lipgloss.JoinVertical(lipgloss.Center, rs...)
}

// TODO: calculate the number of items to show, set v.wb and v.we accordingly
func (v *VerticalViewport) Resize(h int) {

}
