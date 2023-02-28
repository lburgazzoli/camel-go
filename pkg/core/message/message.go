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

	ce.NewEvent()
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
	annotations map[string]string
	content     interface{}
}

func (m *defaultMessage) Fail(err error) {
	m.err = err
}

func (m *defaultMessage) Error() error {
	return nil
}

func (m *defaultMessage) SetAnnotation(key string, val string) {
	if m.annotations == nil {
		m.annotations = make(map[string]string)
	}

	m.annotations[key] = val
}

func (m *defaultMessage) Annotation(key string) (string, bool) {
	if m.annotations == nil {
		return "", false
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
