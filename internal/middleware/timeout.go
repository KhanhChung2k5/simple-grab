package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout middleware
// If the request takes longer than the timeout, the request is aborted and an error is returned
// The request is aborted and an error is returned if the context is done
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			timeout,
		)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
