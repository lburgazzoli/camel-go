package message

import (
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/pkg/errors"
)

func ContentAsBytes(msg camel.Message) ([]byte, error) {
	var content []byte

	_, err := msg.Context().TypeConverter().Convert(msg.Content(), &content)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert message content to []byte")
	}

	return content, nil
}
