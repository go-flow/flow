package flow

import "time"

// RequestLogger returns a middleware that logs all requests on attached router
//
// By default it will log a unique "request_id", the HTTP Method of the request,
// the path that was requested, the duration (time) it took to process the
// request, the size of the response (and the "human" size), and the status
// code of the response.
func RequestLogger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()

		if ct := c.ContentType(); ct != "" {
			c.LogField("content_type", ct)
		}
		c.LogFields(map[string]interface{}{
			"status":    c.Response.Status(),
			"method":    c.Request.Method,
			"path":      c.Request.URL.String(),
			"client_ip": c.ClientIP(),
			"duration":  time.Since(start).String(),
			"size":      c.Response.Size(),
		})
		c.Logger().Info("request-logger")
	}
}
