package flow

// Invoker is implemented by any value that has Invoke method, which invokes provided constructor function
// and returns value created by given constructor. The Invoke function is called for Dependency injection functionalities within flow.
type Invoker interface {
	Invoke(constructor interface{}) (interface{}, error)
}
