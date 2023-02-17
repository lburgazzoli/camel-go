package core

import (
	"github.com/lburgazzoli/camel-go/pkg/core/message"
)

// NewMessage returns a new Message.
func NewMessage() message.Message {
	return message.New()
}
