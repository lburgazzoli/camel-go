package uuid

import "github.com/rs/xid"

func New() string {
	return xid.New().String()
}
