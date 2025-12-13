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

type ResetColor struct{}

func (c ResetColor) Col() string {
	return "\033[0m"
}

type GreenColor struct{}

func (c GreenColor) Col() string {
	return "\033[32m"
}

type BlueColor struct{}

func (c BlueColor) Col() string {
	return "\033[34m"
}

type MagentaColor struct{}

func (c MagentaColor) Col() string {
	return "\033[35m"
}

type BigGreenColor struct{}

func (c BigGreenColor) Col() string {
	return "\033[1;92m"
}

type TealColor struct{}

func (c TealColor) Col() string {
	return "\033[36m"

}

type LightGrayColor struct{}

func (c LightGrayColor) Col() string {
	return "\033[0;37m"
}

type DarkGrayColor struct{}

func (c DarkGrayColor) Col() string {
	return "\033[1;30m"
}

type YellowColor struct{}

func (c YellowColor) Col() string {
	return "\033[38;5;214;3m"
}

type HighlightColor struct{}

func (c HighlightColor) Col() string {
	return "\033[38;5;73m"
	//return "\033[35m"
}

/*

var colors = []string{
	"\033[31m",   // Red
	"\033[32m",   // Green
	"\033[33m",   // Yellow
	"\033[34m",   // Blue
	"\033[35m",   // Magenta
	"\033[36m",   // Cyan
	"\033[37m",   // White
	"\033[1;35m", //light purple
	"\033[1;92m", //Big green
}

func Colors() []string {
	return colors
}
*/
