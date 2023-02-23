package message

import (
	"time"

	"github.com/lburgazzoli/camel-go/pkg/api"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

// New returns a new Message.
func New() (api.Message, error) {
	m := defaultMessage{
		EventContext: &ce.EventContextV1{},
		annotations:  nil,
	}

	if err := m.SetID(uuid.New()); err != nil {
		return nil, err
	}
	if err := m.SetTime(time.Now()); err != nil {
		return nil, err
	}

	return &m, nil
}

type defaultMessage struct {
	ce.EventContext

	err         error
	annotations map[string]interface{}
	content     interface{}
}

func (m *defaultMessage) Fail(err error) {
	m.err = err
}

func (m *defaultMessage) Error() error {
	return nil
}

// SetAnnotation ---
func (m *defaultMessage) SetAnnotation(key string, val interface{}) {
	if m.annotations == nil {
		m.annotations = make(map[string]interface{})
	}

	m.annotations[key] = val
}

// Annotation ---
func (m *defaultMessage) Annotation(key string) (interface{}, bool) {
	if m.annotations == nil {
		return nil, false
	}

	r, ok := m.annotations[key]

	return r, ok
}

func (m *defaultMessage) Content() interface{} {
	return m.content
}

func (m *defaultMessage) SetContent(content interface{}) {
	m.content = content
}
