package flow

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// PanicRecovery returns a middleware that recovers from any panics and writes a 500 if there was one.
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

				httprequest, _ := httputil.DumpRequest(c.Request, false)

				l := c.Logger().WithFields(map[string]interface{}{
					"request": string(httprequest),
					"err":     err,
				})
				l.Error("panic-recovery")

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.ServeError(http.StatusInternalServerError, []byte(err.(error).Error()))
				}
			}
		}()
		c.Next()
	}
}
