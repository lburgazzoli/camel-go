package context

import (
	"io"
)

type Context interface {
	ID() string
	Registry() Registry
	LoadRoutes(in io.Reader) error
}
