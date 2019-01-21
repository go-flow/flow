package flow

// Controller defines set of actions needed in application controller
//
// Controller is used on app#RegisterController
type Controller interface {

	// Routes returns list of controller routes
	Routes() *Router
}
