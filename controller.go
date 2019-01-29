package flow

// ControllerRouter allows controller to define custom controller routing
//
// Controller is used on app#RegisterController
type ControllerRouter interface {

	// Routes returns controller action routing
	Routes() *Router
}

// ControllerPrefixer allows to customize Controller prefix
// for controller action routing
type ControllerPrefixer interface {
	Prefix() string
}

// ControllerIniter allows to init controller
//
type ControllerIniter interface {
	Init(app *App)
}
