package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func PlaceOverlayCenter(bg, fg string, opts ...lipgloss.WhitespaceOption) string {
	return PlaceOverlay(lipgloss.Center, lipgloss.Center, bg, fg, opts...)
}

// PlaceOverlay places fg on top of bg.
func PlaceOverlay(xpos, ypos lipgloss.Position, bg, fg string, opts ...lipgloss.WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	x := int(xpos*lipgloss.Position(bgWidth) - (lipgloss.Position(fgWidth) / 2))
	y := int(ypos*lipgloss.Position(bgHeight) - (lipgloss.Position(fgHeight) / 2))

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			pos = ansi.StringWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(strings.Repeat(" ", x-pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += ansi.StringWidth(fgLine)

		right := ansi.TruncateLeft(bgLine, pos, "")
		bgWidth = ansi.StringWidth(bgLine)
		rightWidth := ansi.StringWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(strings.Repeat(" ", bgWidth-rightWidth-pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) ([]string, int) {
	lines := strings.Split(s, "\n")
	widest := 0
	for _, l := range lines {
		w := ansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
