package flow

// Provider is implemented by any value that has Provide method.
// The Provider instance is used for constructor functions where Dependency injection machanism is decoupled
// Provide method is used to provide constructor value with injected dependecies from given injector
type Provider interface {
	Provide(injector Injector) (interface{}, error)
}

type valueProvider struct {
	c interface{}
}

func (vp *valueProvider) Provide(injector Injector) (interface{}, error) {
	return injector.Inject(vp.c)
}

// NewProvider creates Dependency Injection value provider for given constructor
func NewProvider(constructor interface{}) Provider {
	return &valueProvider{
		c: constructor,
	}
}
