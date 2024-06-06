// //go:build components_http || components_all

package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/lburgazzoli/camel-go/pkg/components"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/pkg/errors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
)

type Producer struct {
	components.DefaultProducer

	endpoint *Endpoint
	tc       api.TypeConverter

	once     sync.Once
	htclient *http.Client
}

func (p *Producer) Endpoint() api.Endpoint {
	return p.endpoint
}

func (p *Producer) Start(context.Context) error {
	return nil
}

func (p *Producer) Stop(context.Context) error {
	return nil
}

func (p *Producer) Receive(ac actor.Context) {
	switch msg := ac.Message().(type) {
	case *actor.Started:
		_ = p.Start(context.Background())
	case *actor.Stopping:
		_ = p.Stop(context.Background())
	case api.Message:
		p.publish(context.Background(), msg)

		// TODO: handle
		if msg.Error() != nil {
			panic(msg.Error())
		}

		ac.Request(ac.Parent(), msg)
	}
}

func (p *Producer) publish(ctx context.Context, msg api.Message) {
	var body []byte

	if msg.Content() != nil {
		body = make([]byte, 0)

		_, err := p.tc.Convert(msg.Content(), &body)
		if err != nil {
			msg.SetError(errors.Wrap(err, "error converting content to []byte"))
			return
		}
	} else {
		body = []byte{}
	}

	req, err := http.NewRequestWithContext(ctx, p.endpoint.config.Method, p.endpoint.config.URI, bytes.NewReader(body))
	if err != nil {
		msg.SetError(err)
		return
	}

	if err := msg.EachHeader(func(k string, v any) error {
		var hval []byte

		_, err := p.tc.Convert(v, &hval)
		if err != nil {
			return err
		}

		req.Header.Set(k, string(hval))

		return nil
	}); err != nil {
		msg.SetError(err)
		return
	}

	res, err := p.client().Do(req)
	if err != nil {
		msg.SetError(err)
		return
	}

	if res.Body != nil {
		defer func() {
			if err := res.Body.Close(); err != nil {
				// TODO: use multierr
				msg.SetError(err)
			}
		}()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			msg.SetError(err)
		}

		msg.SetContent(b)
		msg.SetAttribute(AttributeStatusCode, res.StatusCode)
		msg.SetAttribute(AttributeStatusMessage, res.Status)
		msg.SetAttribute(AttributeProto, res.Proto)
		msg.SetAttribute(AttributeContentLength, res.ContentLength)

		for name, headers := range res.Header {
			val := any(headers)

			if len(headers) == 1 {
				val = headers[0]
			}

			msg.SetHeader(name, val)
		}

		if val := res.Header.Get("Content-Type"); val != "" {
			msg.SetContentType(val)
		}
	}
}

func (p *Producer) client() *http.Client {
	p.once.Do(func() {
		p.htclient = cleanhttp.DefaultClient()
	})

	return p.htclient
}
