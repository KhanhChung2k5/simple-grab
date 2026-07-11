package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/gin-gonic/gin"
)



//Timeout middleware
//If the request takes longer than the timeout, the request is aborted and an error is returned
//The request is aborted and an error is returned if the context is done
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Create a new context with the timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		
		//Create a new channel to signal that the request is finished
		finished := make(chan struct{})
		//Start a new goroutine to handle the request
		go func() {
			//Handle the request
			c.Next()
			close(finished)
		}()

		//Select the finished channel or the context done channel
		select {
		//If the request is finished, do nothing
		case <-finished:
		case <-ctx.Done():
			response.Error(c, errors.New("request timeout"))
			c.Abort()
			return
		}
	}
}
