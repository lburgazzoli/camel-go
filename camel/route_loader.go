package camel

// ==========================
//
// Route Loader
//
// ==========================

// RouteLoader --
type RouteLoader interface {
	Load(data []byte) ([]Definition, error)
}
