package console

//takes an array of strings and pads the right of the shorter strings
//so they are all the same length.

func NormalisePadding(items []string) []string {
	length := MaxLen(items)
	var paddedItems []string
	for _, item := range items {
		fitted := FitString(item, length)
		paddedItems = append(paddedItems, fitted)
	}
	return paddedItems
}

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

func FitString(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	var str = s
	leng := len(s)
	for i := leng; i <= length; i++ {
		str = str + " "
	}
	return str
}
