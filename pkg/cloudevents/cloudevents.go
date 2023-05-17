package cloudevents

import (
	"context"
	"encoding/json"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/mimetype"
	"github.com/pkg/errors"
)

func AsJSON(_ context.Context, msg camel.Message) ([]byte, error) {
	payload := CloudEventJSON{
		SpecVersion:       "1.0",
		ID:                msg.ID(),
		Type:              msg.Type(),
		Source:            msg.Source(),
		Subject:           msg.Subject(),
		Time:              msg.Time().String(),
		DataContentType:   msg.ContentType(),
		DataContentSchema: msg.ContentSchema(),
	}

	if msg.Content() != nil {
		if payload.DataContentType == "" {
			payload.DataContentType = mimetype.ApplicationOctetStream
		}

		// TODO: the conversion function should take into account the content type

		switch payload.DataContentType {
		case mimetype.ApplicationOctetStream:
			content, err := message.ContentAsBytes(msg)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert message content to []byte")
			}

			EncodeDataContentAsBase64(&payload, content)
		case mimetype.ApplicationJSON:
			content, err := message.ContentAsBytes(msg)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert message content to []byte")
			}

			payload.Data = content
		case mimetype.ApplicationXML:
			content, err := message.ContentAsBytes(msg)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert message content to []byte")
			}

			payload.Data = content
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert message content to []byte")
	}

	return bytes, err
}
