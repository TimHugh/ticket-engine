package mock

import (
	"fmt"
	"strings"
)

type Logger struct {
	Messages []string
}

func (l *Logger) Printf(msg string, params ...interface{}) {
	l.Messages = append(l.Messages, fmt.Sprintf(msg, params...))
}

func (l *Logger) Contains(msg string) bool {
	for _, entry := range l.Messages {
		if strings.Contains(entry, msg) {
			return true
		}
	}

	return false
}

func (l *Logger) Out() string {
	var msg string
	for _, entry := range l.Messages {
		msg += entry + "\n"
	}
	return msg
}
