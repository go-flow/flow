package flow

// Starter interface used when http application is served
// Start method is invoked if module implements the interface
type Starter interface {
	Start() error
}

// Stopper interface used when http application is stopped
// Stop method is invoked during shutdown process if module implements the interface
type Stopper interface {
	Stop()
}

// Bootstrap creates flow Module instance for given factory object
func Bootstrap(factory interface{}) (*Module, error) {
	module, err := NewModule(factory, nil)
	return module, err
}
