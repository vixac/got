package console

type GotTheme struct {
}

func (t *GotTheme) ColorFor(token Token) Color {
	switch token.Name() {
	case "primary":
		return GreenColor{}
	case "brand":
		return RedColor{}
	case "secondary":
		return Resetcolor{}
	case "complete":
		return MagentaColor{}
	}
	return Resetcolor{}
}
