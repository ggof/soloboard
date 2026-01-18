package viewport

import (
	"iter"
)

type Viewport struct {
	I        int
	beg      int
	end      int
	max      int
	size     int
	itemSize int
	len      int
}

func New(itemSize int) Viewport {
	return Viewport{itemSize: itemSize}
}

func (v *Viewport) Prev() {
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

func (v *Viewport) Next() {
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

func (v *Viewport) SetSize(size int) {
	v.size = size
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

	v.max = max(1, v.size/(v.itemSize))

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
