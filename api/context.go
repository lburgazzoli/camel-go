package api

// HasContext --
type HasContext interface {
	Context() Context
}

// ContextAware --
type ContextAware interface {
	HasContext

	SetContext(context Context)
}

// Context --
type Context interface {
	Service

	// Registry
	Registry() LoadingRegistry

	// Type conversion
	AddTypeConverter(converter TypeConverter)
	//TypeConverters() []TypeConverter
	TypeConverter() TypeConverter

	// Routes
	AddRoute(route *Route)
	//Routes() []*Route

	// Services
	AddService(service Service) bool
	//Service() []Service
}
