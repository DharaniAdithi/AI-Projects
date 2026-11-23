package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"booking_system/booking_service/internal/config"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/internal/repository"
	"booking_system/booking_service/internal/services"
	"booking_system/booking_service/pkg/database"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SetupRoutes(router *mux.Router, db *database.Db, logger *loggers.Logger, cfg *config.Config) {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(db, logger)

	// Initialize services
	bookingService := services.NewBookingService(bookingRepo, logger)

	// Initialize handlers
	bookingHandler := NewBookingHandler(bookingService, logger)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Booking routes
	api.HandleFunc("/booking", bookingHandler.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods("GET")
	api.HandleFunc("/bookings/{uuid}", bookingHandler.GetBookingByUUID).Methods("GET")
	api.HandleFunc("/bookings/{uuid}", bookingHandler.PatchBooking).Methods("PATCH")
	// Add this route to the booking routes section:
	api.HandleFunc("/bookings/{uuid}/cancel", bookingHandler.CancelBooking).Methods("PATCH")
	api.HandleFunc("/bookings/{uuid}/recover", bookingHandler.RecoverBooking).Methods("PATCH")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "booking-service"}`))
	}).Methods(http.MethodGet)

	// Serve OpenAPI spec
	router.HandleFunc("/booking_service/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		openAPIPath := filepath.Join(cfg.ProjectRoot, "docs", "openapi.yaml")
		yamlContent, err := os.ReadFile(openAPIPath)
		if err != nil {
			fmt.Println("Error reading openAPI.yaml:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		updatedContent := strings.ReplaceAll(string(yamlContent), "{BASE_URL}", cfg.BaseUrl)
		w.Header().Set("Content-Type", "application/x-yaml")
		if _, writeErr := w.Write([]byte(updatedContent)); writeErr != nil {
			logger.Info("Error writing response", zap.Error(writeErr))
		}
	}).Methods(http.MethodGet)

	// Serve Swagger UI
	swaggerUIPath := http.Dir(filepath.Join(cfg.ProjectRoot, "docs", "swaggerui"))
	fs := http.FileServer(swaggerUIPath)
	router.PathPrefix("/booking/docs/").Handler(http.StripPrefix("/booking/docs/", fs))

	logger.Info("Routes configured successfully")

	seedInitialData(bookingService)

	logger.Info("Data seeded successfully")
}

func seedInitialData(bookingService services.BookingServiceInterface) {
	ctx := context.Background()

	users := []models.User{
		{UserName: "John Doe", Email: "john@example.com", Phone: "9876543210", IsActive: true},
		{UserName: "Jane Smith", Email: "jane@example.com", Phone: "9876543211", IsActive: true},
	}

	rooms := []models.Room{
		{RoomName: "Conference Room A", Capacity: 10, PricePerDay: 5000, IsActive: true},
		{RoomName: "Conference Room B", Capacity: 20, PricePerDay: 8000, IsActive: true},
		{RoomName: "Meeting Room", Capacity: 5, PricePerDay: 2000, IsActive: true},
	}

	bookingService.SeedUsersAndRooms(ctx, users, rooms)
}
