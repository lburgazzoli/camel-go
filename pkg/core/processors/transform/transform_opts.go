package transform

import "github.com/lburgazzoli/camel-go/pkg/core/language"

type OptionFn func(*Transform)

func WithLanguage(lang language.Language) OptionFn {
	return func(in *Transform) {
		in.Language = lang
	}
}
