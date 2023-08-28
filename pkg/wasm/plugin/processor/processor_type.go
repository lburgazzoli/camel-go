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
	Data          []byte            `json:"data,omitempty"`
	Headers       map[string][]byte `json:"headers,omitempty"`
	Attributes    map[string][]byte `json:"attributes,omitempty"`
}

type Evaluation struct {
	Result bool `json:"result,omitempty"`
}
