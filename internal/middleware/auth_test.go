package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeAuthProvider struct {
	claims   *service.Claims
	user     *model.User
	parseErr error
	userErr  error
}

func (f *fakeAuthProvider) ParseToken(
	token string,
) (*service.Claims, error) {
	if f.parseErr != nil {
		return nil, f.parseErr
	}
	return f.claims, nil
}

func (f *fakeAuthProvider) GetMe(
	ctx context.Context,
	userID string,
) (*model.User, error) {
	if f.userErr != nil {
		return nil, f.userErr
	}
	return f.user, nil
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validUser := &model.User{ID: "user-1", Role: model.RoleRider}

	tests := []struct {
		name       string
		header     string
		provider   *fakeAuthProvider
		wantStatus int
	}{
		{name: "missing authorization header", provider: &fakeAuthProvider{}, wantStatus: http.StatusUnauthorized},
		{name: "invalid authorization scheme", header: "Basic token", provider: &fakeAuthProvider{}, wantStatus: http.StatusUnauthorized},
		{name: "missing bearer token", header: "Bearer ", provider: &fakeAuthProvider{}, wantStatus: http.StatusUnauthorized},
		{name: "invalid token", header: "Bearer invalid", provider: &fakeAuthProvider{parseErr: errors.New("bad token")}, wantStatus: http.StatusUnauthorized},
		{name: "user no longer exists", header: "Bearer valid", provider: &fakeAuthProvider{claims: &service.Claims{UserID: "user-1"}, userErr: errors.New("not found")}, wantStatus: http.StatusUnauthorized},
		{name: "valid token", header: "Bearer valid", provider: &fakeAuthProvider{claims: &service.Claims{UserID: "user-1"}, user: validUser}, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/protected", AuthMiddleware(tt.provider), func(c *gin.Context) {
				if c.GetString("user_id") != validUser.ID || c.GetString("role") != validUser.Role {
					t.Fatal("auth middleware did not set user context")
				}
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			if res.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d; body=%s", tt.wantStatus, res.Code, res.Body.String())
			}
		})
	}
}
