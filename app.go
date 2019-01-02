package flow

import "sync"

// App -
type App struct {
	Options
	router *Router
	pool   sync.Pool
}
