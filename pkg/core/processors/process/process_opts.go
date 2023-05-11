package process

type OptionFn func(*Process)

func WithRef(ref string) OptionFn {
	return func(in *Process) {
		in.Ref = ref
	}
}
