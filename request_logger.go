package flow

import (
	"time"

	"github.com/rs/xid"
)

// RequestLogger returns a middleware that logs all requests on attached router
//
// By default it will log a unique "request_id", the HTTP Method of the request,
// the path that was requested, the duration (time) it took to process the
// request, the size of the response (and the "human" size), and the status
// code of the response.
func RequestLogger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()

		// check if request ID exists in headers
		requestID := c.RequestID()

		if requestID == "" {
			// generate new RequestID
			guid := xid.New()
			requestID = guid.String()
			// add requestID to header
			c.Request.Header.Add("X-Request-ID", requestID)
		}

		//execute next handler in chain
		c.Next()

		l := c.Logger().WithFields(map[string]interface{}{
			"request_id": requestID,
			"status":     c.Response.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.String(),
			"client_ip":  c.ClientIP(),
			"duration":   time.Since(start).String(),
			"human_size": byteCountDecimal(int64(c.Response.Size())),
			"size":       c.Response.Size(),
		})
		l.Info("request-logger")
	}
}
