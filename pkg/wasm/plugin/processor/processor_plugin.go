//go:build tinygo.wasm

package processor

import (
	"context"
	"io"
	"os"
)

type Processor func(context.Context, *Message) (*Message, error)

var processor Processor

func RegisterProcessors(p Processor) {
	processor = p
}

//
//export process
func _process(size uint32) uint64 {
	b := make([]byte, size)

	_, err := io.ReadAtLeast(os.Stdin, b, int(size))
	if err != nil {
		return 0
	}

	req := Message{}
	if err := req.UnmarshalVT(b); err != nil {
		return 0
	}
	response, err := processor(context.Background(), &req)
	if err != nil {
		n, err := os.Stdout.WriteString(err.Error())
		if err != nil {
			return 0
		}

		// Indicate that this is the error string by setting the 32-th bit, assuming that
		// no data exceeds 31-bit size (2 GiB).
		return uint64(n) | 1<<31
	}

	b, err = response.MarshalVT()
	if err != nil {
		return 0
	}

	n, err := os.Stdout.Write(b)
	if err != nil {
		return 0
	}

	return uint64(n)
}