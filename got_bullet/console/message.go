package console

import (
	"fmt"
	"unicode/utf8"
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
		length += utf8.RuneCountInString(m.Message)
	}
	return MessageGroup{
		Messages: messages,
		TextLen:  length,
	}
}

func (m Message) InColor() string {
	return fmt.Sprintf("%s%s%s", m.Color, m.Message, ResetColor{}.Col())
}
