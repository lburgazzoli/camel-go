package wasm

import (
	"context"
	"errors"
	"fmt"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	wzapi "github.com/tetratelabs/wazero/api"
)

type Function struct {
	module *Module
	fn     wzapi.Function
}

//nolint:mnd,gosec
func (p *Function) invoke(ctx context.Context, msg camel.Message) error {
	ctx = context.WithValue(ctx, contextKeyModule, p.module)
	ctx = context.WithValue(ctx, contextKeyMessage, msg)

	ret, err := p.fn.Call(ctx)
	if err != nil {
		return err
	}

	if ret[0] != 0 {
		resPtr := uint32(ret[0] >> 32)
		resLen := uint32(ret[0])

		switch uint8(resLen >> 28) {
		case 0xF:
			// error
			size := resLen & 0x0FFFFFFF

			errText, ok := p.module.Memory().Read(resPtr, size)
			if !ok {
				err = fmt.Errorf(
					"memory.Read(%d, %d) out of range of memory size %d",
					resPtr,
					size,
					p.module.Memory().Size(),
				)
			} else {
				err = errors.New(string(errText))
			}

			msg.SetError(err)
		case 0x1:
			// true
			// TODO: maybe better to have some better result rather than using the error
			//       as a signal for predicate true/false
			return ErrPredicateMatches
		case 0x2:
			// false
			// TODO: maybe better to have some better result rather than using the error
			//       as a signal for predicate true/false
			return ErrPredicateDoesNotMatch
		}
	}

	return nil
}
