package flow

type Provider interface {
	Provide(injector Injector) (interface{}, error)
}

type instanceProvider struct {
	constructor interface{}
}

func (ip *instanceProvider) Provide(injector Injector) (interface{}, error) {
	return injector.Provide(ip.constructor)
}

func NewProvider(constructor interface{}) Provider {
	return &instanceProvider{
		constructor: constructor,
	}
}
