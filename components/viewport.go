package components

import (
	"iter"
)

type Viewport struct {
	I   int
	beg int
	end int
	max int
	h   int
	iht int
	len int
}

func NewViewport(length int, itemHeight int) *Viewport {
	return &Viewport{len: length, iht: itemHeight}
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
	if v.I < v.len-1 {
		v.I++

		for v.I >= v.end {
			v.end++
			for v.end-v.beg > v.max {
				v.beg++
			}
		}
	}
}

func (v *Viewport) SetSize(h int) {
	v.h = h
	v.updateViewportSize()
}

func (v *Viewport) SetLen(length int) {
	v.len = length
	v.updateViewportSize()
}

func (v *Viewport) updateViewportSize() {
	if v.len == 0 {
		return
	}

	v.max = max(1, v.h/(v.iht+1)) // count some margin, show at least 1 item

	v.beg = v.I
	v.end = v.I
	for i := range v.max {
		if i%2 == 0 {
			if v.beg > 0 {
				v.beg--
			} else if v.end < v.len {
				v.end++
			}
		} else {
			if v.end < v.len {
				v.end++
			} else if v.beg > 0 {
				v.beg--
			}
		}
	}
}

func (v Viewport) Window() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := v.beg; i < v.end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}
