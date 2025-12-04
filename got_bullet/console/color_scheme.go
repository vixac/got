package console

type Theme interface {
	ColorFor(token Token) Color
}

//make a color scheme

type RedColor struct {
}

func (c RedColor) Col() string {
	return "\033[31m"
}

type Resetcolor struct{}

func (c Resetcolor) Col() string {
	return resetColor
}

type GreenColor struct{}

func (c GreenColor) Col() string {
	return "\033[32m"
}

type BlueColor struct{}

func (c BlueColor) Color() string {
	return "\033[34m"
}

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
