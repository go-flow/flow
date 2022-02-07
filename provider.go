package flow

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
