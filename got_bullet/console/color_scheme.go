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

type AmberColor struct{}

func (c AmberColor) Col() string {
	return "\033[38;5;185m"
}

type HighlightColor struct{}

func (c HighlightColor) Col() string {
	return "\033[38;5;73m"
}
