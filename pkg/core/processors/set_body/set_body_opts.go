package setbody

import "github.com/lburgazzoli/camel-go/pkg/core/language"

type OptionFn func(*SetBody)

func WithLanguage(lang language.Language) OptionFn {
	return func(in *SetBody) {
		in.Language = lang
	}
}
