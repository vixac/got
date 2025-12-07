package console

import (
	"unicode/utf8"
)

// max length of input items
func MaxLen(items []string) int {
	var max = 0
	for _, item := range items {
		l := utf8.RuneCountInString(item)
		if l > max {
			max = l
		}
	}
	return max
}

func FitString(s string, length int, pad string) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= length {
		return s[:length]
	}
	var str = s
	for i := runeCount; i < length; i++ {
		str = str + pad
	}
	return str
}
