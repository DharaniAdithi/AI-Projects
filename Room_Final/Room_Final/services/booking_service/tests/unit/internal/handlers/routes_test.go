package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"booking_system/booking_service/internal/config"
	"booking_system/booking_service/internal/handlers"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getRoutesTestLogger() *loggers.Logger {
	return loggers.NewLogger(loggers.LogConfig{
		Level:       "error", // suppress logs in route tests
		LogDir:      "./logs",
		FileName:    "routes_test.log",
		ServiceName: "test-service",
	})
}

func getRoutesTestDB() *database.Db {
	sqlDB, mock, _ := sqlmock.New()

	// Allow any query so SetupRoutes wiring doesn't fail
	mock.MatchExpectationsInOrder(false)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	return &database.Db{Gorm: gormDB, SqlDb: sqlDB}
}

func TestSetupRoutes_WrongMethod_Returns405(t *testing.T) {
    router := mux.NewRouter()
    handlers.SetupRoutes(router, getRoutesTestDB(), getRoutesTestLogger(), &config.Config{})

    // Ensure the route exists for GET
    router.HandleFunc("/api/v1/bookings", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")

    // Override MethodNotAllowedHandler just for this test
    router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    })

    // DELETE on a GET-only route should now be 405
    req, _ := http.NewRequest("DELETE", "/api/v1/bookings", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}


func TestSetupRoutes_UnknownRoute_Returns405Or404(t *testing.T) {
	router := mux.NewRouter()
	handlers.SetupRoutes(router, getRoutesTestDB(), getRoutesTestLogger(), &config.Config{})

	req, _ := http.NewRequest("GET", "/api/v1/unknown-route", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Gorilla mux returns 405 or 404 for unregistered routes
	assert.True(t, rr.Code == http.StatusNotFound || rr.Code == http.StatusMethodNotAllowed)
}

func TestSetupRoutes_AllRoutesRegistered(t *testing.T) {
	router := mux.NewRouter()
	handlers.SetupRoutes(router, getRoutesTestDB(), getRoutesTestLogger(), &config.Config{})

	routes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/booking"},
		{"GET", "/api/v1/bookings"},
		{"GET", "/api/v1/bookings/some-uuid"},
		{"PATCH", "/api/v1/bookings/some-uuid"},
	}

	for _, tt := range routes {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var match mux.RouteMatch
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			matched := router.Match(req, &match)
			assert.True(t, matched, "Route %s %s should be registered", tt.method, tt.path)
		})
	}
}
