package flow

// Injector is implemented by any value that has Inject method, which invokes provided constructor function
// and returns value created by given constructor. The Inject function is called for Dependency injection functionalities within flow.
type Injector interface {
	Inject(constructor interface{}) (interface{}, error)
}
