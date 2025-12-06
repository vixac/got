package console

import "unicode/utf8"

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
	if utf8.RuneCountInString(s) > length {
		return s[:length]
	}
	var str = s
	leng := utf8.RuneCountInString(s)
	for i := leng; i <= length; i++ {
		str = str + pad
	}
	return str
}
