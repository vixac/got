package console

import (
	"fmt"
)

type Message struct {
	Message string
	Color   string
}

func (m Message) InColor() string {
	return fmt.Sprintf("%s%s%s", m.Color, m.Message, resetColor)
}
