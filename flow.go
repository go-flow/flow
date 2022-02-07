package flow

// Bootstrap creates flow Module instance for given factory object
func Bootstrap(factory interface{}) (*Module, error) {
	module, err := NewModule(factory, nil)
	return module, err
}
