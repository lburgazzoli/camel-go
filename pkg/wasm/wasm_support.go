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
