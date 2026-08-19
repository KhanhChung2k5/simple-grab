package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/gin-gonic/gin"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		userRole       string
		requiredRole   string
		expectedStatus int
	}{
		{
			name:           "rider can access rider route",
			userRole:       model.RoleRider,
			requiredRole:   model.RoleRider,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "driver can access driver route",
			userRole:       model.RoleDriver,
			requiredRole:   model.RoleDriver,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "rider cannot access driver route",
			userRole:       model.RoleRider,
			requiredRole:   model.RoleDriver,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "missing role is forbidden",
			userRole:       "",
			requiredRole:   model.RoleRider,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.GET(
				"/test",
				func(c *gin.Context) {
					if tt.userRole != "" {
						c.Set("role", tt.userRole)
					}
					c.Next()
				},
				RequireRole(tt.requiredRole),
				func(c *gin.Context) {
					c.Status(http.StatusOK)
				},
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			if res.Code != tt.expectedStatus {
				t.Fatalf(
					"expected %d, got %d",
					tt.expectedStatus,
					res.Code,
				)
			}
		})
	}
}