package processor

import (
	"time"
)

type Message struct {
	ID            string            `json:"id,omitempty"`
	Source        string            `json:"source,omitempty"`
	Type          string            `json:"type,omitempty"`
	Subject       string            `json:"subject,omitempty"`
	ContentType   string            `json:"content_type,omitempty"`
	ContentSchema string            `json:"content_schema,omitempty"`
	Time          time.Time         `json:"time,omitempty"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	Data          []byte            `json:"data,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

func (x *Message) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Message) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Message) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *Message) GetContentSchema() string {
	if x != nil {
		return x.ContentSchema
	}
	return ""
}

func (x *Message) GetTime() time.Time {
	if x != nil {
		return x.Time
	}
	return time.Time{}
}

func (x *Message) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Message) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}
