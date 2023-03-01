package interop

import karmem "karmem.org/golang"

func DecodeMessage(data []byte) Message {
	reader := karmem.NewReader(data)
	decoded := NewMessageViewer(reader, 0)

	m := Message{
		ID:            decoded.ID(reader),
		Source:        decoded.Source(reader),
		Type:          decoded.Type(reader),
		Subject:       decoded.Subject(reader),
		ContentType:   decoded.ContentType(reader),
		ContentSchema: decoded.ContentSchema(reader),
		Time:          decoded.Time(),
		Content:       decoded.Content(reader),
	}

	decoded.Annotations(reader)

	for _, a := range decoded.Annotations(reader) {
		m.Annotations = append(m.Annotations, Annotation{
			Key: a.Key(reader),
			Val: a.Val(reader),
		})
	}

	return m
}
