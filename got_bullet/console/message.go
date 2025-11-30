package console

import (
	"fmt"
)

const resetColor = "\033[0m"
const BIGreen = "\033[1;92m"

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
	"\033[37m", // White
}

func Colors() []string {
	return colors
}

type Message struct {
	Message string
	Color   string
}

func (m Message) InColor() string {
	return fmt.Sprintf("%s[%s] %s", m.Color, m.Message, resetColor)
}
