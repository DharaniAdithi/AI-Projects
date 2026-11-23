package handlers

import (
    "net/http"        // ADD
    "path/filepath"   // ADD

    "booking_system/booking_service/internal/config"
    "booking_system/booking_service/internal/loggers"
    "booking_system/booking_service/internal/repository"
    "booking_system/booking_service/internal/services"
    "booking_system/booking_service/pkg/database"

    "github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, db *database.Db, logger *loggers.Logger, cfg *config.Config) {
    bookingRepo := repository.NewBookingRepository(db, logger)
    bookingService := services.NewBookingService(bookingRepo, logger)
    bookingHandler := NewBookingHandler(bookingService, logger)

    api := router.PathPrefix("/api/v1").Subrouter()

    api.HandleFunc("/booking", bookingHandler.CreateBooking).Methods("POST")
    api.HandleFunc("/bookings", bookingHandler.GetAllBookings).Methods("GET")
    api.HandleFunc("/bookings/{uuid}", bookingHandler.GetBookingByUUID).Methods("GET")
    api.HandleFunc("/bookings/{uuid}", bookingHandler.PatchBooking).Methods("PATCH")

    // --- Swagger UI (no impact on API routes above) ---
    docsDir := filepath.Join(cfg.ProjectRoot, "docs")
    router.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(docsDir, "openapi.yaml"))
    })
    router.PathPrefix("/swagger/").Handler(
        http.StripPrefix("/swagger/", http.FileServer(http.Dir(filepath.Join(docsDir, "swaggerui")))),
    )

    logger.Info("Routes configured successfully")
}
