package flow

// ControllerRouter interface allows controllers to define their routing logic
type ControllerRouter interface {
	Routes(*Router)
}
