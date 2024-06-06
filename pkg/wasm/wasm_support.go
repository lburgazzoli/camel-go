package wasm

type contextKeyType string

const (
	noPointer = 0
	noSize    = 0

	contextKeyModule  contextKeyType = "_module"
	contextKeyMessage contextKeyType = "_message"

	allocFunctionNAme   = "alloc"
	deallocFunctionNAme = "dealloc"
)

func toPtrSize(ptr uint64, size uint64) uint64 {
	return ptr<<32 | size
}
