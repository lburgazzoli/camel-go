package wasm

import (
	"context"
	"fmt"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	wzapi "github.com/tetratelabs/wazero/api"
)

func BuiltInFunctions() []HostFunction {
	return []HostFunction{
		{
			Name:    "message_get_id",
			Fn:      getMessageID,
			Params:  []wzapi.ValueType{},
			Results: []wzapi.ValueType{wzapi.ValueTypeI64},
		},

		{
			Name:    "message_get_content",
			Fn:      getMessageContent,
			Params:  []wzapi.ValueType{},
			Results: []wzapi.ValueType{wzapi.ValueTypeI64},
		},

		{
			Name:    "message_set_content",
			Fn:      setMessageContent,
			Params:  []wzapi.ValueType{wzapi.ValueTypeI32, wzapi.ValueTypeI32},
			Results: []wzapi.ValueType{},
		},
	}
}

func getMessageID(ctx context.Context, mod *Module, msg camel.Message, stack []uint64) error {
	val := msg.ID()

	ptr, size, err := mod.write(ctx, []byte(val))
	if err != nil {
		return err
	}

	stack[0] = toPtrSize(ptr, size)

	return nil
}

func getMessageContent(ctx context.Context, mod *Module, msg camel.Message, stack []uint64) error {
	var content []byte

	if msg.Content() != nil {
		content = make([]byte, 0)

		_, err := msg.Context().TypeConverter().Convert(msg.Content(), &content)
		if err != nil {
			msg.SetError(fmt.Errorf("error converting content to []byte: %w", err))
		}
	} else {
		content = []byte{}
	}

	ptr, size, err := mod.write(ctx, content)
	if err != nil {
		return err
	}

	stack[0] = toPtrSize(ptr, size)

	return nil
}

func setMessageContent(_ context.Context, mod *Module, msg camel.Message, stack []uint64) error {
	offset := wzapi.DecodeU32(stack[0])
	size := wzapi.DecodeU32(stack[1])

	data, ok := mod.Memory().Read(offset, size)
	if !ok {
		panic(
			fmt.Errorf(
				"memory.Read(%d, %d) out of range of memory size %d",
				offset,
				size,
				mod.module.Memory().Size(),
			),
		)
	}

	dst := make([]byte, len(data))
	copy(dst, data)

	msg.SetContent(dst)

	return nil
}
