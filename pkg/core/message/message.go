package message

import (
	"maps"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

// New returns a new Message.
func New(ctx api.Context) api.Message {
	m := defaultMessage{
		ctx:        ctx,
		attributes: make(map[string]any),
		meta: defaultMessageMeta{
			ID:   uuid.New(),
			Time: time.Now(),
		},
	}

	return &m
}

type defaultMessageMeta struct {
	ID            string
	Source        string
	Type          string
	Subject       string
	ContentType   string
	ContentSchema string
	Time          time.Time
}

type defaultMessage struct {
	ctx api.Context
	err error

	meta       defaultMessageMeta
	headers    map[string]any
	attributes map[string]any
	content    interface{}
}

func (m *defaultMessage) Context() api.Context {
	return m.ctx
}

func (m *defaultMessage) ID() string {
	return m.meta.ID
}

func (m *defaultMessage) Time() time.Time {
	return m.meta.Time
}

func (m *defaultMessage) Type() string {
	return m.meta.Type
}

func (m *defaultMessage) SetType(v string) {
	m.meta.Type = v
}

func (m *defaultMessage) Source() string {
	return m.meta.Source
}

func (m *defaultMessage) SetSource(v string) {
	m.meta.Source = v
}

func (m *defaultMessage) Subject() string {
	return m.meta.Subject
}

func (m *defaultMessage) SetSubject(v string) {
	m.meta.Subject = v
}

func (m *defaultMessage) ContentSchema() string {
	return m.meta.ContentSchema
}

func (m *defaultMessage) SetContentSchema(v string) {
	m.meta.ContentSchema = v
}

func (m *defaultMessage) ContentType() string {
	return m.meta.ContentType
}

func (m *defaultMessage) SetContentType(v string) {
	m.meta.ContentType = v
}

//
// Errors
//

func (m *defaultMessage) SetError(err error) {
	m.err = err
}

func (m *defaultMessage) Error() error {
	return m.err
}

//
// Content
//

func (m *defaultMessage) Content() interface{} {
	return m.content
}

func (m *defaultMessage) SetContent(content interface{}) {
	m.content = content
}

//
// Attributes
//

func (m *defaultMessage) Attributes() map[string]any {
	answer := make(map[string]any)

	maps.Copy(answer, m.attributes)

	return answer
}

func (m *defaultMessage) SetAttributes(attributes map[string]any) {
	m.attributes = make(map[string]any)

	maps.Copy(attributes, m.attributes)
}

func (m *defaultMessage) Attribute(key string) (any, bool) {
	if m.attributes == nil {
		return nil, false
	}

	r, ok := m.attributes[key]

	return r, ok
}

func (m *defaultMessage) SetAttribute(key string, val any) {
	if m.attributes == nil {
		m.attributes = make(map[string]any)
	}

	m.attributes[key] = val
}

func (m *defaultMessage) EachAttribute(fn func(string, any) error) error {
	for k, v := range m.attributes {
		if err := fn(k, v); err != nil {
			return err
		}
	}

	return nil
}

//
// Headers
//

func (m *defaultMessage) Headers() map[string]any {
	answer := make(map[string]any)

	maps.Copy(answer, m.headers)

	return answer
}

func (m *defaultMessage) SetHeaders(headers map[string]any) {
	m.headers = make(map[string]any)

	maps.Copy(headers, m.headers)
}

func (m *defaultMessage) Header(key string) (any, bool) {
	if m.headers == nil {
		return nil, false
	}

	r, ok := m.headers[key]

	return r, ok
}

func (m *defaultMessage) SetHeader(key string, val any) {
	if m.headers == nil {
		m.headers = make(map[string]any)
	}

	m.headers[key] = val
}

func (m *defaultMessage) EachHeader(fn func(string, any) error) error {
	for k, v := range m.headers {
		if err := fn(k, v); err != nil {
			return err
		}
	}

	return nil
}

//
// Clone
//

func (m *defaultMessage) CopyTo(_ api.Message) error {
	// TODO implement me.
	panic("implement me")
}
