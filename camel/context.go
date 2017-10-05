package camel

// Context --
type Context struct {
	name            string
	registryLoaders []RegistryLoader
	components      map[string]Component
}

// ==========================
//
// Initialize a camel context
//
// ==========================

// NewContext --
func NewContext() *Context {
	return &Context{
		name:            "camel",
		registryLoaders: make([]RegistryLoader, 0),
		components:      make(map[string]Component),
	}
}

// NewContextWithName --
func NewContextWithName(name string) *Context {
	return &Context{
		name:            name,
		registryLoaders: make([]RegistryLoader, 0),
		components:      make(map[string]Component),
	}
}

// ==========================
//
//
//
// ==========================

// AddRegistryLoader --
func (context *Context) AddRegistryLoader(loader RegistryLoader) {
	context.registryLoaders = append(context.registryLoaders, loader)
}

// AddTypeConverter --
func (context *Context) AddTypeConverter(converter TypeConverter) {
}

// AddComponent --
func (context *Context) AddComponent(name string, component Component) {
	context.components[name] = component
	context.components[name].SetContext(context)
}

// Component --
func (context *Context) Component(name string) (Component, error) {
	component, found := context.components[name]

	// check if the component has already been created or added to the context
	// component list
	if !found {
		for _, loader := range context.registryLoaders {
			component, err := loader.Load(name)

			if err != nil {
				return nil, err
			}

			if component == nil {
				continue
			}

			if _, ok := component.(Component); !ok {
				// not a component
				continue
			}

			if component != nil {
				break
			}
		}

		if component != nil {
			context.AddComponent(name, component)
		}
	}

	return component, nil
}

// ==========================
//
// Lyfecycle
//
// ==========================

// Start --
func (context *Context) Start() {
}

// Stop --
func (context *Context) Stop() {
}
