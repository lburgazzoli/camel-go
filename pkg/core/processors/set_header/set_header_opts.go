package setheader

import "github.com/lburgazzoli/camel-go/pkg/core/language"

type OptionFn func(*SetHeader)

func WithLanguage(lang language.Language) OptionFn {
	return func(in *SetHeader) {
		in.Language = lang
	}
}

func WithName(name string) OptionFn {
	return func(in *SetHeader) {
		in.Name = name
	}
}
