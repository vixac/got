package console

import (
	"fmt"
)

type Messenger interface {
	Print(message Message)
	PrintInLine(line []Message)
	Error(message Message)
}
type Printer struct {
}

func (p Printer) PrintInLine(line []Message) {
	for _, m := range line {
		fmt.Printf("%s", m.InColor())
	}
	fmt.Printf("\n")
}

func (p Printer) Print(message Message) {
	fmt.Println(message.InColor())
}

func (p Printer) Error(message Message) {
	fmt.Println("Error:" + message.Message)
}
