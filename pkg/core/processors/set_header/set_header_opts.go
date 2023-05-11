package setheader

type OptionFn func(*SetHeader)

func WithLanguage(lang Language) OptionFn {
	return func(in *SetHeader) {
		in.Language = lang
	}
}

func WithName(name string) OptionFn {
	return func(in *SetHeader) {
		in.Name = name
	}
}
