package serdes

import (
	"time"

	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/wasm/interop"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"

	"sync"

	karmem "karmem.org/golang"
)

var writerPool = sync.Pool{
	New: func() any {
		return karmem.NewWriter(1024)
	},
}

func EncodeMessage(message camel.Message) ([]byte, error) {
	writer, ok := writerPool.Get().(*karmem.Writer)
	if !ok {
		panic("no writable pool")
	}
	defer func() {
		writer.Reset()
		writerPool.Put(writer)
	}()

	content := interop.Message{
		ID:            message.GetID(),
		Source:        message.GetSource(),
		Type:          message.GetType(),
		Subject:       message.GetSubject(),
		ContentType:   message.GetDataContentType(),
		ContentSchema: message.GetDataSchema(),
		Time:          message.GetTime().UnixMilli(),
	}

	message.ForEachAnnotation(func(k string, v string) {
		content.Annotations = append(content.Annotations, interop.Pair{
			Key: k,
			Val: v,
		})
	})

	if message.Content() != nil {
		switch d := message.Content().(type) {
		case []byte:
			content.Content = d
		case string:
			content.Content = []byte(d)
		default:
			return nil, camelerrors.InternalErrorf("unsupported type %v", message.Content())
		}
	}

	if _, err := content.WriteAsRoot(writer); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func DecodeMessage(encoded []byte) (camel.Message, error) {

	msg, err := message.New()
	if err != nil {
		return nil, err
	}

	if err := DecodeMessageTo(encoded, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func DecodeMessageTo(encoded []byte, msg camel.Message) error {
	reader := karmem.NewReader(encoded)
	decoded := interop.NewMessageViewer(reader, 0)

	_ = msg.SetID(decoded.ID(reader))
	_ = msg.SetSource(decoded.Source(reader))
	_ = msg.SetType(decoded.Type(reader))
	_ = msg.SetSubject(decoded.Subject(reader))
	_ = msg.SetDataContentType(decoded.ContentType(reader))
	_ = msg.SetDataSchema(decoded.ContentSchema(reader))
	_ = msg.SetTime(time.UnixMilli(decoded.Time()))

	msg.SetContent(decoded.Content(reader))

	for _, a := range decoded.Annotations(reader) {
		msg.SetAnnotation(
			a.Key(reader),
			a.Val(reader),
		)
	}

	return nil
}
