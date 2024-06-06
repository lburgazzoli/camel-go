package cloudevents

import "encoding/base64"

const (
	Spec1_0 = "1.0"

	ExtensionPartitionKey = "partitionkey"
	ExtensionSequence     = "sequence"
)

//nolint:tagliatelle
type CloudEventJSON struct {
	SpecVersion       string `json:"specversion"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	Source            string `json:"source"`
	Subject           string `json:"subject,omitempty"`
	Time              string `json:"time"`
	Data              []byte `json:"data,omitempty"`
	DataBase64        []byte `json:"data_base64,omitempty"`
	DataContentType   string `json:"datacontenttype,omitempty"`
	DataContentSchema string `json:"datacontentschmea,omitempty"`
}

func EncodeDataContentAsBase64(ce *CloudEventJSON, content []byte) {
	ce.DataBase64 = make([]byte, base64.StdEncoding.EncodedLen(len(content)))
	base64.StdEncoding.Encode(ce.DataBase64, content)
}
