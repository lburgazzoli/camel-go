package language

import (
	"github.com/lburgazzoli/camel-go/pkg/core/language/constant"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/lburgazzoli/camel-go/pkg/core/language/mustache"
)

type OptionFn func(language *Language)

func WithJq(value *jq.Jq) OptionFn {
	return func(in *Language) {
		in.Jq = value
	}
}
func WithJqExpression(value string) OptionFn {
	return func(in *Language) {
		in.Jq = &jq.Jq{
			Definition: jq.Definition{
				Expression: value,
			},
		}
	}
}

func WithConstant(value *constant.Constant) OptionFn {
	return func(in *Language) {
		in.Constant = value
	}
}
func WithConstantExpression(value string) OptionFn {
	return func(in *Language) {
		in.Constant = &constant.Constant{
			Value: value,
		}
	}
}

func WithMustache(value *mustache.Mustache) OptionFn {
	return func(in *Language) {
		in.Mustache = value
	}
}
func WithMustacheExpression(value string) OptionFn {
	return func(in *Language) {
		in.Mustache = &mustache.Mustache{
			Template: value,
		}
	}
}
