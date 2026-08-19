package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"gorm.io/gorm"
)

type apiResponse struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
}

func TestIntegrationRegisterRideComplete(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to a dedicated PostgreSQL test database")
	}
	parsed, err := url.Parse(databaseURL)
	if err != nil || !strings.Contains(strings.ToLower(strings.Trim(parsed.Path, "/")), "test") {
		t.Fatal("TEST_DATABASE_URL must point to a database whose name contains 'test'")
	}

	ctx := context.Background()
	db, err := repository.ConnectDB(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql database: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	applyTestMigrations(t, db)

	userRepo := repository.NewUserRepository(db)
	rideRepo := repository.NewRideRepository(db)
	driverRepo := repository.NewDriverRepository(db)
	authService := service.NewAuthService(userRepo, "phase-6-test-secret", 1)
	router := NewRouter(
		authService,
		handler.NewAuthHandler(authService),
		handler.NewRideHandler(service.NewRideService(rideRepo, driverRepo)),
		handler.NewDriverHandler(service.NewDriverService(driverRepo)),
		handler.NewHealthHandler(),
	)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	riderEmail := "phase6-rider-" + suffix + "@example.com"
	driverEmail := "phase6-driver-" + suffix + "@example.com"
	riderID := register(t, router, riderEmail, "rider")
	driverID := register(t, router, driverEmail, "driver")
	t.Cleanup(func() {
		db.Exec("DELETE FROM rides WHERE rider_id = ? OR driver_id = ?", riderID, driverID)
		db.Exec("DELETE FROM users WHERE id IN ?", []string{riderID, driverID})
	})

	riderToken := login(t, router, strings.ToUpper(riderEmail))
	driverToken := login(t, router, driverEmail)

	ride := request(t, router, http.MethodPost, "/api/v1/rides", riderToken, map[string]any{
		"pickup_lat": 10.762622, "pickup_lng": 106.660172,
		"dropoff_lat": 10.776889, "dropoff_lng": 106.700806,
	}, http.StatusOK)
	var rideData struct{ ID, Status string }
	decodeData(t, ride, &rideData)
	if rideData.ID == "" || rideData.Status != "pending" {
		t.Fatalf("unexpected created ride: %+v", rideData)
	}

	request(t, router, http.MethodPost, "/api/v1/rides", riderToken, map[string]any{
		"pickup_lat": 10.7, "pickup_lng": 106.6, "dropoff_lat": 10.8, "dropoff_lng": 106.7,
	}, http.StatusConflict)
	request(t, router, http.MethodPatch, "/api/v1/drivers/me/location", driverToken, map[string]any{
		"latitude": 91, "longitude": 106.7,
	}, http.StatusBadRequest)
	request(t, router, http.MethodPatch, "/api/v1/drivers/me/location", driverToken, map[string]any{
		"latitude": 10.77, "longitude": 106.7,
	}, http.StatusOK)
	request(t, router, http.MethodPatch, "/api/v1/drivers/me/online", driverToken, map[string]any{
		"is_online": true,
	}, http.StatusOK)
	request(t, router, http.MethodPost, "/api/v1/rides/"+rideData.ID+"/accept", driverToken, nil, http.StatusOK)
	request(t, router, http.MethodPatch, "/api/v1/rides/"+rideData.ID+"/status", driverToken, map[string]any{
		"status": "completed",
	}, http.StatusConflict)
	request(t, router, http.MethodPatch, "/api/v1/rides/"+rideData.ID+"/status", driverToken, map[string]any{
		"status": "in_progress",
	}, http.StatusOK)
	completed := request(t, router, http.MethodPatch, "/api/v1/rides/"+rideData.ID+"/status", driverToken, map[string]any{
		"status": "completed",
	}, http.StatusOK)
	decodeData(t, completed, &rideData)
	if rideData.Status != "completed" {
		t.Fatalf("expected completed ride, got %q", rideData.Status)
	}
	request(t, router, http.MethodGet, "/api/v1/rides/not-a-uuid", riderToken, nil, http.StatusBadRequest)
	request(t, router, http.MethodGet, "/api/v1/rides/00000000-0000-4000-8000-000000000000", riderToken, nil, http.StatusNotFound)
}

func applyTestMigrations(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, name := range []string{"001_init.up.sql", "002_active_ride_constraints.up.sql"} {
		path := filepath.Join("..", "..", "migrations", name)
		sql, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if err := db.Exec(string(sql)).Error; err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func register(t *testing.T, router http.Handler, email, role string) string {
	t.Helper()
	res := request(t, router, http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"email": email, "password": "grab123", "role": role,
	}, http.StatusOK)
	var data struct{ ID string }
	decodeData(t, res, &data)
	if data.ID == "" {
		t.Fatal("register response has no user id")
	}
	return data.ID
}

func login(t *testing.T, router http.Handler, email string) string {
	t.Helper()
	res := request(t, router, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email": email, "password": "grab123",
	}, http.StatusOK)
	var data struct {
		AccessToken string `json:"access_token"`
	}
	decodeData(t, res, &data)
	if data.AccessToken == "" {
		t.Fatal("login response has no access token")
	}
	return data.AccessToken
}

func request(t *testing.T, router http.Handler, method, path, token string, payload any, wantStatus int) apiResponse {
	t.Helper()
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Fatalf("encode request: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &body)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	var decoded apiResponse
	if err := json.Unmarshal(res.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode response (%d): %v; body=%s", res.Code, err, res.Body.String())
	}
	if res.Code != wantStatus {
		t.Fatalf("%s %s: expected %d, got %d; error=%q", method, path, wantStatus, res.Code, decoded.Error)
	}
	return decoded
}

func decodeData(t *testing.T, response apiResponse, target any) {
	t.Helper()
	if err := json.Unmarshal(response.Data, target); err != nil {
		t.Fatalf("decode response data: %v", err)
	}
}
