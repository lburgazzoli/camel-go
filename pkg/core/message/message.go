package message

import (
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

// New returns a new Message.
func New() Message {
	m := Message{
		Event:       ce.NewEvent(),
		annotations: nil,
	}

	m.SetID(uuid.New())
	m.SetTime(time.Now())

	return m
}

type Message struct {
	ce.Event

	annotations map[string]interface{}
}

// Annotation ---
// TODO: this may be a private method to store contextual info for routing purposes
func (m *Message) Annotation(key string, val interface{}) {
	if m.annotations == nil {
		m.annotations = make(map[string]interface{})
	}

	m.annotations[key] = val
}

// GetAnnotation ---
// TODO: this may be a private method to store contextual info for routing purposes
func (m *Message) GetAnnotation(key string) (interface{}, bool) {
	if m.annotations == nil {
		return nil, false
	}

	r, ok := m.annotations[key]

	return r, ok
}
