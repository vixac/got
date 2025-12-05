package console

type GotTheme struct {
}

func (t *GotTheme) ColorFor(token Token) Color {
	switch token.Name() {
	case "primary":
		return GreenColor{}
	case "brand":
		return BigGreenColor{}
	case "secondary":
		return ResetColor{}
	case "complete":
		return MagentaColor{}
	case "gid":
		return BlueColor{}
	}
	return ResetColor{}
}
