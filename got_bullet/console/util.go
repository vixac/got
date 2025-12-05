package console

// max length of input items
func MaxLen(items []string) int {
	var max = 0
	for _, item := range items {
		l := len(item)
		if l > max {
			max = l
		}
	}
	return max
}

func FitString(s string, length int, pad string) string {
	if len(s) > length {
		return s[:length]
	}
	var str = s
	leng := len(s)
	for i := leng; i <= length; i++ {
		str = str + pad
	}
	return str
}
