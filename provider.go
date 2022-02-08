package flow

type ImportsProvider interface {
	// ProvideImports returns list of instance providers for module dependecies
	// This method is used to register all module dependecies
	// eg. logging, db connection,....
	// all dependecies that are provided in this method
	// will be available to all modules imported by the factory
	ProvideImports() []Provider
}

type ExportsProvider interface {

	// ProvideExports returns list of instance providers for
	// functionalities that module will export.
	// Exported functionalities will be available to other modules that
	// import module created by the Factory
	ProvideExports() []Provider
}

type ModulesProvider interface {

	// ProvideModules returns list of instance providers
	// for modules that current module depends on
	ProvideModules() []Provider
}

type RouterProvider interface {

	// ProvideRouters returns list of instance providers for module routers.
	// Module routers are used for http routing
	ProvideRouters() []Provider
}

// Provider is implemented by any value that has Provide method.
// The Provider instance is used for constructor functions where Dependency injection machanism is decoupled
// Provide method is used to provide constructor value with injected dependecies from given injector
type Provider interface {
	Provide(in Invoker) (interface{}, error)
}

type valueProvider struct {
	c interface{}
}

func (vp *valueProvider) Provide(in Invoker) (interface{}, error) {
	return in.Invoke(vp.c)
}

// NewValueProvider creates Dependency Injection value provider for given constructor
func NewValueProvider(c interface{}) Provider {
	return &valueProvider{
		c: c,
	}
}
