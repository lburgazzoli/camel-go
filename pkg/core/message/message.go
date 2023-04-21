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

func (m *defaultMessage) SetError(err error) {
	m.err = err
}

func (m *defaultMessage) Error() error {
	return m.err
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

func (m *defaultMessage) Annotations() map[string]string {
	answer := make(map[string]string)

	for k, v := range m.annotations {
		answer[k] = v
	}

	return answer

}
func (m *defaultMessage) SetAnnotations(annotations map[string]string) {
	m.annotations = make(map[string]string)

	for k, v := range annotations {
		m.annotations[k] = v
	}
}

func (m *defaultMessage) ForEachAnnotation(fn func(string, string)) {
	for k, v := range m.annotations {
		fn(k, v)
	}
}

func (m *defaultMessage) Content() interface{} {
	return m.content
}

func (m *defaultMessage) SetContent(content interface{}) {
	m.content = content
}

func (m *defaultMessage) CopyTo(message api.Message) error {

	if err := message.SetID(m.GetID()); err != nil {
		return err
	}
	if err := message.SetSource(m.GetSource()); err != nil {
		return err
	}
	if err := message.SetType(m.GetType()); err != nil {
		return err
	}
	if err := message.SetSubject(m.GetSubject()); err != nil {
		return err
	}
	if err := message.SetDataContentType(m.GetDataContentType()); err != nil {
		return err
	}
	if err := message.SetDataSchema(m.GetDataSchema()); err != nil {
		return err
	}
	if err := message.SetTime(m.GetTime()); err != nil {
		return err
	}

	message.SetContent(m.Content())

	message.SetAnnotations(m.Annotations())
	message.SetError(m.Error())

	for k, v := range m.GetExtensions() {
		if err := message.SetExtension(k, v); err != nil {
			return err
		}
	}

	return nil
}
