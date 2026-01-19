package utils

func EllipsisBeg(text string, w int) string {
	if len(text) > w {
		tbs := []byte(text[len(text)-w:])
		copy(tbs[0:3], "...")
		text = string(tbs)
	}
	return text
}

func EllipsisEnd(text string, w int) string {
	if len(text) > w {
		tbs := []byte(text[:w])
		copy(tbs[w-3:], "...")
		text = string(tbs)
	}
	return text
}
