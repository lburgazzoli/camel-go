//go:build tinygo.wasm

package processor

import (
	"context"
	"encoding/json"
	"io"
	"os"
)

type Processor func(context.Context, *Message) (*Message, error)

var processor Processor

func RegisterProcessors(p Processor) {
	processor = p
}

//export process
func _process(size uint32) uint64 {
	b := make([]byte, size)

	_, err := io.ReadAtLeast(os.Stdin, b, int(size))
	if err != nil {
		return 0
	}

	req := Message{}
	if err := json.Unmarshal(b, &req); err != nil {
		return 0
	}
	res, err := processor(context.Background(), &req)
	if err != nil {
		n, err := os.Stdout.WriteString(err.Error())
		if err != nil {
			return 0
		}

		// Indicate that this is the error string
		return (uint64(1) << uint64(32)) | uint64(n)
	}

	b, err = json.Marshal(res)
	if err != nil {
		return 0
	}

	n, err := os.Stdout.Write(b)
	if err != nil {
		return 0
	}

	return uint64(n)
}
