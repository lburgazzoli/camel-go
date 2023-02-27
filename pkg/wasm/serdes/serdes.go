package serdes

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"time"

	karmem "karmem.org/golang"
	"sync"
)

var writerPool = sync.Pool{
	New: func() any {
		return karmem.NewWriter(1024)
	},
}

func Encode(message api.Message) ([]byte, error) {
	writer := writerPool.Get().(*karmem.Writer)
	defer func() {
		writer.Reset()
		writerPool.Put(writer)
	}()

	content := Message{
		ID:            message.GetID(),
		Source:        message.GetSource(),
		Type:          message.GetType(),
		Subject:       message.GetSubject(),
		ContentType:   message.GetDataContentType(),
		ContentSchema: message.GetDataSchema(),
		Time:          message.GetTime().UnixMilli(),
	}

	if _, err := content.WriteAsRoot(writer); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func Decode(encoded []byte) (api.Message, error) {
	reader := karmem.NewReader(encoded)
	decoded := NewMessageViewer(reader, 0)

	msg, err := message.New()
	if err != nil {
		return nil, err
	}

	_ = msg.SetID(decoded.ID(reader))
	_ = msg.SetSource(decoded.Source(reader))
	_ = msg.SetType(decoded.Type(reader))
	_ = msg.SetSubject(decoded.Subject(reader))
	_ = msg.SetDataContentType(decoded.ContentType(reader))
	_ = msg.SetDataSchema(decoded.ContentSchema(reader))
	_ = msg.SetTime(time.UnixMilli(decoded.Time()))

	return msg, nil
}
