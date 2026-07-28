package middleware

import (
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
	"strings"
)

//Error constants (invalid token)
var (
	ErrInvalidToken = errors.New("invalid token")
)

//AuthMiddleware is a middleware for the auth service
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get the token from the header
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" { //No token provided
			//Return an error if the token is empty 
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		//Split the token
		tokenParts := strings.SplitN(tokenStr, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || tokenParts[1] == "" { //Bearer <token> (sai format)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		//Get the token
		token := tokenParts[1] //<token>
		//Parse the token
		claims, err := authService.ParseToken(token)
		if err != nil { //Invalid token
			//Return an error if the token is invalid
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		//Get the user by id
		user, err := authService.GetMe(c.Request.Context(), claims.UserID)
		if err != nil { //User not found
			//Return an error if the user is not found
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		//Set the user in the context
		c.Set("user", user) //user is the user object
		c.Set("user_id", user.ID) //user_id is the user id
		c.Set("role", user.Role) //role is the user role
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
