package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger is a middleware function that logs the details of incoming requests.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		agent := c.Request.UserAgent()

		var body string
		if method == "POST" {
			bodyBytes, err := c.GetRawData()
			if err != nil {
				logger.Error("Error reading request body", zap.Error(err))
			} else {
				body = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

		}

		c.Next()

		status := c.Writer.Status()

		logger.Info("Request Info",
			zap.String("path", path),
			zap.String("method", method),
			zap.String("ip", ip),
			zap.String("agent", agent),
			zap.Int("status", status),
			zap.String("body", body),
		)

	}
}
