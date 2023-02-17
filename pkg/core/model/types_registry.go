package model

type TypeFactory func() interface{}

var Types = make(map[string]TypeFactory)
