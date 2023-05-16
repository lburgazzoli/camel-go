package cloudevents

type CloudEventJSON struct {
	SpecVersion       string `json:"specversion"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	Source            string `json:"source"`
	Subject           string `json:"subject,omitempty"`
	Time              string `json:"time"`
	Data              []byte `json:"data,omitempty"`
	DataContentType   string `json:"datacontenttype,omitempty"`
	DataContentSchema string `json:"datacontentschmea,omitempty"`
}
