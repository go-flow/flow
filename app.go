package flow

import "sync"

// App -
type App struct {
	config *Config
	router *Router
	pool   sync.Pool
}
