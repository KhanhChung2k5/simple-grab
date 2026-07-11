package middleware

import (
	"log"
	"errors"

	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler handles unprocessed request errors.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		log.Printf("request error: %v", err)
		response.InternalServerError(c, errors.New("internal server error"))
	}
}
