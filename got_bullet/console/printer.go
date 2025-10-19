package console

import (
	"fmt"
)

type Messenger interface {
	Print(message Message)
	Error(message Message)
}
type Printer struct {
}

func (p Printer) Print(message Message) {
	fmt.Println(message.Message)
}

func (p Printer) Error(message Message) {
	fmt.Println("Error:" + message.Message)
}
