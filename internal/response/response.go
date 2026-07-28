package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Body struct {
	Data any `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}
//Status 200
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Data: data})
}

//Status 500
func Error(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Body{Error: err.Error()})
}

//Status 400
func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Body{Error: err.Error()})
}

//Status 404
func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, Body{Error: err.Error()})
}

//Status 500
func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Body{Error: err.Error()})
}

//Status 401
func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Body{Error: err.Error()})
}

//Status 403
func Forbidden(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, Body{Error: err.Error()})
}

// Status 409
func Conflict(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, Body{Error: err.Error()})
}

//Status 429
func TooManyRequests(c *gin.Context, err error) {
	c.JSON(http.StatusTooManyRequests, Body{Error: err.Error()})
}

//Status 502
func BadGateway(c *gin.Context, err error) {
	c.JSON(http.StatusBadGateway, Body{Error: err.Error()})
}

//Status 503
func ServiceUnavailable(c *gin.Context, err error) {
	c.JSON(http.StatusServiceUnavailable, Body{Error: err.Error()})
}

//Status 504
func GatewayTimeout(c *gin.Context, err error) {
	c.JSON(http.StatusGatewayTimeout, Body{Error: err.Error()})
}



