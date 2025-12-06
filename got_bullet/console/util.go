package console

import (
	"fmt"
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
	if length == runeCount {
		//		fmt.Printf("VX: This string wants to be zero '%s', %d \n", s, length)
		//		return s
	}
	if runeCount >= length {
		fmt.Printf("VX: shortcutting")
		return s[:length]
	}
	var str = s
	for i := runeCount; i < length; i++ {
		str = str + pad
	}
	return str
}
