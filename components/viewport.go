package components

import (
	"github.com/charmbracelet/lipgloss"
)

type List interface {
	RenderItem(i int, selected bool) string
	Len() int
}

type Viewport struct {
	I    int
	beg  int
	end  int
	max  int
	w    int
	h    int
	list List
}

func NewViewport(list List) *Viewport {
	return &Viewport{list: list}
}

func (v *Viewport) Up() {
	if v.I > 0 {
		v.I--

		for v.I < v.beg {
			v.beg--

			for v.end-v.beg > v.max {
				v.end--
			}
		}
	}
}

func (v *Viewport) Down() {
	if v.I < v.list.Len()-1 {
		v.I++

		for v.I >= v.end {
			v.end++

			for v.end-v.beg > v.max {
				v.beg++
			}
		}
	}
}

func (v *Viewport) SetSize(w, h int) {
	v.w, v.h = w, h
	v.updateViewportSize()
}

func (v *Viewport) SetList(list List) {
	v.list = list
	v.updateViewportSize()
}

func (v *Viewport) updateViewportSize() {
	if v.list.Len() == 0 {
		return
	}

	itemHeight := lipgloss.Height(v.list.RenderItem(0, false))
	v.max = max(1, v.h/(itemHeight+1)) // count some margin, show at least 1 item

	v.beg = v.I
	v.end = v.I
	for i := range v.max {
		if i%2 == 0 {
			if v.beg > 0 {
				v.beg--
			} else if v.end < v.list.Len() {
				v.end++
			}
		} else {
			if v.end < v.list.Len() {
				v.end++
			} else if v.beg > 0 {
				v.beg--
			}
		}
	}
}

func (v Viewport) Render() string {
	var items []string
	for i := v.beg; i < v.end; i++ {
		items = append(items, v.list.RenderItem(i, i == v.I))
	}

	return lipgloss.Place(
		v.w,
		v.h,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, items...),
	)
}
