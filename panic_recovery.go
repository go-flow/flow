package flow

import (
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// PanicRecovery returns a middleware that recovers from any panics and serves error response
func PanicRecovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {

				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.ServeError(http.StatusInternalServerError, err.(error))
				}
			}
		}()
		c.Next()
	}
}
