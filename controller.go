package flow

// Controller defines set of actions needed in application controller
//
// Controller is used on app#RegisterController
type Controller interface {

	// Routes returns controller action routing
	Routes() *Router
}

// ControllerPrefixer allows to customize Controller prefix
// for controller action routing
type ControllerPrefixer interface {
	Prefix() string
}
