package console

import (
	"fmt"
)

type Message struct {
	Message string
	Color   string
}

type MessageGroup struct {
	Messages []Message
	TextLen  int
}

// use this with printinline to get length of all messages for the divider
func NewMessageGroup(messages []Message) MessageGroup {
	var length = 0
	for _, m := range messages {
		length += len(m.Message)
	}
	return MessageGroup{
		Messages: messages,
		TextLen:  length,
	}
}

func (m Message) InColor() string {
	return fmt.Sprintf("%s%s%s", m.Color, m.Message, resetColor)
}
