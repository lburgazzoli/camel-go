package processors

type TypeFactory func() interface{}

var Types = make(map[string]TypeFactory)
