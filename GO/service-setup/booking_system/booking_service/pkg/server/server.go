package server

import (
	"net/http"

	"booking_system/booking_service/internal/config"
	"booking_system/booking_service/internal/handlers"
	utils "booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/pkg/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Application struct {
	Logger    *utils.Logger
	Router    *mux.Router
	Config    *config.Config
	DBService database.DBConnector
}

func InitializeApp(env string) (*Application, error) {
	// Load the configuration
	cfg, err := config.LoadConfig("booking-service", env)
	if err != nil {
		return nil, err
	}

	// Initialize logger
	logger := utils.NewLogger(utils.LogConfig{
		Level:       cfg.Log.Level,
		LogDir:      cfg.Log.Directory,
		FileName:    cfg.Log.FileName,
		ServiceName: "booking-service",
	})
	logger.Info("Logger initialized")

	// Initialize database service
	dbService := &database.DBService{}
	logger.Info("Initializing database connection")

	db, err := dbService.EstablishConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}
	logger.Info("Database connection established successfully")

	// Run auto migration
	if err := database.AutoMigrate(db, logger); err != nil {
		logger.Error("Failed to run database migrations", zap.Error(err))
		return nil, err
	}

	// Setup router and routes
	router := mux.NewRouter()
	handlers.SetupRoutes(router, db, logger, cfg)

	return &Application{
		Logger:    logger,
		Router:    router,
		Config:    cfg,
		DBService: dbService,
	}, nil
}

// GetCorsOptions returns the CORS configuration options
func GetCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Access-Token",
			"Device-Identifier",
		},
	}
}

func RunServer(app *Application) error {
	app.Logger.Info(
		"Starting HTTP server",
		zap.String("port", app.Config.Port),
	)

	handler := cors.New(GetCorsOptions()).Handler(app.Router)
	return http.ListenAndServe(":"+app.Config.Port, handler)
}
