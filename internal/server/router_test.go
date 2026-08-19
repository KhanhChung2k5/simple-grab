package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

type routerAuthStub struct {
	user *model.User
}

func (s routerAuthStub) ParseToken(token string) (*service.Claims, error) {
	if token != "valid" || s.user == nil {
		return nil, errors.New("invalid token")
	}
	return &service.Claims{UserID: s.user.ID, Role: s.user.Role}, nil
}

func (s routerAuthStub) GetMe(context.Context, string) (*model.User, error) {
	if s.user == nil {
		return nil, errors.New("user not found")
	}
	return s.user, nil
}

func TestRouterErrorCodeContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	health := handler.NewHealthHandler()

	tests := []struct {
		name       string
		user       *model.User
		header     string
		path       string
		wantStatus int
	}{
		{name: "health is public", path: "/api/v1/health", wantStatus: http.StatusOK},
		{name: "missing token", path: "/api/v1/rider/ping", wantStatus: http.StatusUnauthorized},
		{name: "invalid token", header: "Bearer wrong", path: "/api/v1/rider/ping", wantStatus: http.StatusUnauthorized},
		{name: "wrong role", user: &model.User{ID: "driver-1", Role: model.RoleDriver}, header: "Bearer valid", path: "/api/v1/rider/ping", wantStatus: http.StatusForbidden},
		{name: "correct role", user: &model.User{ID: "rider-1", Role: model.RoleRider}, header: "Bearer valid", path: "/api/v1/rider/ping", wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(routerAuthStub{user: tt.user}, nil, nil, nil, health)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d; body=%s", tt.wantStatus, res.Code, res.Body.String())
			}
		})
	}
}
