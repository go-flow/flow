package flow

// Bootstrap creates flow Module instance for given factory object
func Bootstrap(mf ModuleFactory) (*Module, error) {
	module, err := NewModule(mf, nil)
	return module, err
}
