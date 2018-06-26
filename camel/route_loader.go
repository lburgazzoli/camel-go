package camel

// ==========================
//
// Route Loader
//
// ==========================

// RouteLoader --
type RouteLoader interface {
	Load() ([]Definition, error)
}
