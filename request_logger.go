package flow

import (
	"time"

	"github.com/go-flow/flow/log"

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

		c.Response.Header().Add("X-Request-ID", requestID)

		//c.LogField("request_id", requestID)
		c.LogFields(log.String("request_id", requestID))

		//execute next handler in chain
		c.Next()

		c.LogFields(
			log.Int("status", c.Response.Status()),
			log.String("method", c.Request.Method),
			log.String("path", c.Request.URL.String()),
			log.String("client_ip", c.ClientIP()),
			log.Duration("duration", time.Since(start)),
			log.String("human_size", byteCountDecimal(int64(c.Response.Size()))),
			log.Int("size", c.Response.Size()),
		)
		c.Logger().Info("request-logger")
	}
}
