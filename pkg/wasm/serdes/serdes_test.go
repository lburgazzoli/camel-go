package serdes

import (
	"context"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestSerde(t *testing.T) {
	m, err := message.New()
	assert.Nil(t, err)

	encoded, err := Encode(m)
	assert.Nil(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := Decode(encoded)
	assert.Nil(t, err)

	assert.Equal(t, m.GetID(), decoded.GetID())
	assert.Equal(t, m.GetTime().UnixMilli(), decoded.GetTime().UnixMilli())
}

func TestWASM(t *testing.T) {
	// https://github.com/tetratelabs/wazero/blob/main/examples/allocation/tinygo/greet.go

	ctx := context.Background()

	r := wazero.NewRuntime(ctx)
	defer func() { _ = r.Close(ctx) }()

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	wasmContext, err := os.ReadFile("../../../etc/fn/simple_process.wasm")
	assert.Nil(t, err)

	mod, err := r.Instantiate(ctx, wasmContext)
	assert.Nil(t, err)

	p := mod.ExportedFunction("process")
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	m, err := message.New()
	assert.Nil(t, err)

	encoded, err := Encode(m)
	assert.Nil(t, err)

	results, err := malloc.Call(ctx, uint64(len(encoded)))
	assert.Nil(t, err)

	encodedPtr := results[0]

	defer func() { _, _ = free.Call(ctx, encodedPtr) }()

	assert.True(t, mod.Memory().Write(uint32(encodedPtr), encoded))

	ptrSize, err := p.Call(ctx, encodedPtr, uint64(len(encoded)))
	assert.Nil(t, err)

	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	outPtr := uint32(ptrSize[0] >> 32)
	outSize := uint32(ptrSize[0])

	bytes, ok := mod.Memory().Read(outPtr, outSize)
	assert.True(t, ok)

	decoded, err := Decode(bytes)
	assert.Nil(t, err)
	assert.Equal(t, m.GetID(), decoded.GetID())

	c, ok := decoded.Content().([]byte)
	assert.True(t, ok)
	assert.Equal(t, "hello from wasm", string(c))

}
