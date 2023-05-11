package wasm

type OptionFn func(*Wasm)

func WithRef(ref string) OptionFn {
	return func(in *Wasm) {
		if err := in.UnmarshalText([]byte(ref)); err != nil {
			// TODO: better handling
			panic(err)
		}
	}
}
