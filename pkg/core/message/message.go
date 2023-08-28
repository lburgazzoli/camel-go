package message

import (
	"fmt"
	"maps"
	"time"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

// New returns a new Message.
func New(ctx api.Context) api.Message {
	m := defaultMessage{
		ctx:        ctx,
		attributes: make(map[string]any),
	}

	m.attributes[api.MessageAttributeID] = uuid.New()
	m.attributes[api.MessageAttributeTime] = time.Now()

	return &m
}

type defaultMessage struct {
	ctx        api.Context
	err        error
	headers    map[string]any
	attributes map[string]any
	content    interface{}
}

func (m *defaultMessage) Context() api.Context {
	return m.ctx
}

func (m *defaultMessage) ID() string {
	v, ok := m.attributes[api.MessageAttributeID]
	if !ok {
		panic(camelerrors.InternalErrorf("undefined %s", api.MessageAttributeID))
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeID))
	}

	return answer
}

func (m *defaultMessage) Time() time.Time {
	v, ok := m.attributes[api.MessageAttributeTime]
	if !ok {
		panic(camelerrors.InternalErrorf("undefined %s", api.MessageAttributeTime))
	}

	answer, ok := v.(time.Time)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeTime))
	}

	return answer
}

func (m *defaultMessage) Type() string {
	v, ok := m.attributes[api.MessageAttributeType]
	if !ok {
		return ""
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeType))
	}

	return answer
}

func (m *defaultMessage) SetType(v string) {
	m.attributes[api.MessageAttributeType] = v
}

func (m *defaultMessage) Source() string {
	v, ok := m.attributes[api.MessageAttributeSource]
	if !ok {
		return ""
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeSource))
	}

	return answer
}

func (m *defaultMessage) SetSource(v string) {
	m.attributes[api.MessageAttributeSource] = v
}

func (m *defaultMessage) Subject() string {
	v, ok := m.attributes[api.MessageAttributeSubject]
	if !ok {
		return ""
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeSubject))
	}

	return answer
}

func (m *defaultMessage) SetSubject(v string) {
	m.attributes[api.MessageAttributeSubject] = v
}

func (m *defaultMessage) ContentSchema() string {
	v, ok := m.attributes[api.MessageAttributeContentSchema]
	if !ok {
		return ""
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeContentSchema))
	}

	return answer
}

func (m *defaultMessage) SetContentSchema(v string) {
	m.attributes[api.MessageAttributeContentSchema] = v
}

func (m *defaultMessage) ContentType() string {
	v, ok := m.attributes[api.MessageAttributeContentType]
	if !ok {
		return ""
	}

	answer, ok := v.(string)
	if !ok {
		panic(camelerrors.InternalErrorf("wrong type for %s", api.MessageAttributeContentType))
	}

	return answer
}

func (m *defaultMessage) SetContentType(v string) {
	m.attributes[api.MessageAttributeContentType] = v
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

func (m *defaultMessage) validateAttribute(key string) error {
	switch key {
	case api.MessageAttributeID:
		return fmt.Errorf("attempt to set a reserved attribute: %s", key)
	case api.MessageAttributeTime:
		return fmt.Errorf("attempt to set a reserved attribute: %s", key)
	default:
		return nil
	}
}

func (m *defaultMessage) Attributes() map[string]any {
	answer := make(map[string]any)

	maps.Copy(answer, m.attributes)

	return answer
}

func (m *defaultMessage) SetAttributes(attributes map[string]any) error {
	m.attributes = make(map[string]any)
	m.attributes[api.MessageAttributeID] = m.ID()
	m.attributes[api.MessageAttributeTime] = m.Time()

	for k, v := range attributes {
		if err := m.validateAttribute(k); err != nil {
			return err
		}

		m.attributes[k] = v
	}

	return nil
}

func (m *defaultMessage) Attribute(key string) (any, bool) {
	r, ok := m.attributes[key]

	return r, ok
}

func (m *defaultMessage) SetAttribute(key string, val any) error {
	if err := m.validateAttribute(key); err != nil {
		return err
	}

	m.attributes[key] = val

	return nil
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
