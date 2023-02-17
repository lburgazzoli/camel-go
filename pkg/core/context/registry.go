package context

type Registry interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}
