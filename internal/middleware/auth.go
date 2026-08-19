package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

// Error constants (invalid token)
var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthProvider interface {
	ParseToken(token string) (*service.Claims, error)
	GetMe(
		ctx context.Context,
		userID string,
	) (*model.User, error)
}

// AuthMiddleware is a middleware for the auth service
func AuthMiddleware(authService AuthProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get the token from the header
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" { //No token provided
			//Abort 401 if the token is empty
			response.AbortUnauthorized(c, ErrInvalidToken)
			return
		}

		//Split the token
		tokenParts := strings.SplitN(tokenStr, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || tokenParts[1] == "" { //Bearer <token> (sai format)
			//Abort 401 if the token is invalid
			response.AbortUnauthorized(c, ErrInvalidToken)
			return
		}
		//Get the token
		token := tokenParts[1] //<token>
		//Parse the token
		claims, err := authService.ParseToken(token)
		if err != nil { //Invalid token
			//Return an error if the token is invalid
			response.AbortUnauthorized(c, ErrInvalidToken)
			return
		}
		//Get the user by id
		user, err := authService.GetMe(c.Request.Context(), claims.UserID)
		if err != nil { //User not found
			//Return an error if the user is not found
			response.AbortUnauthorized(c, ErrInvalidToken)
			return
		}
		//Set the user in the context
		c.Set("user", user)       //user is the user object
		c.Set("user_id", user.ID) //user_id is the user id
		c.Set("role", user.Role)  //role is the user role
		//Continue
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		userRole := c.GetString("role")
		if _, ok := allowed[userRole]; !ok {
			response.AbortForbidden(c, errors.New("forbidden"))
			return
		}
		c.Next()
	}
}
