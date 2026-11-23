package database

import (
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"

	"go.uber.org/zap"
)

func AutoMigrate(db *Db, logger *loggers.Logger) error {
	logger.Info("Starting auto migration for database tables")

	err := db.Gorm.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Booking{},
	)

	if err != nil {
		logger.Error("Failed to run auto migration", zap.Error(err))
		return err
	}

	logger.Info("Auto migration completed successfully")
	return nil
}
