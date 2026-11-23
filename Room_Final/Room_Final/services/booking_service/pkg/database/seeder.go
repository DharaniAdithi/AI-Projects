package database

import (
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"

	"go.uber.org/zap"
)

func SeedInitialData(db *Db, logger *loggers.Logger) {
	logger.Info("Seeding initial data")

	users := []models.User{
		{UserName: "John Doe", Email: "john@example.com", Phone: "9876543210", IsActive: true},
		{UserName: "Jane Smith", Email: "jane@example.com", Phone: "9876543211", IsActive: true},
	}

	rooms := []models.Room{
		{RoomName: "Conference Room A", Capacity: 10, IsActive: true},
		{RoomName: "Conference Room B", Capacity: 20, IsActive: true},
		{RoomName: "Meeting Room", Capacity: 5, IsActive: true},
	}

	for _, user := range users {
		if err := db.Gorm.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			logger.Error("Failed to seed user", zap.Error(err), zap.String("email", user.Email))
		}
	}

	for _, room := range rooms {
		if err := db.Gorm.FirstOrCreate(&room, models.Room{RoomName: room.RoomName}).Error; err != nil {
			logger.Error("Failed to seed room", zap.Error(err), zap.String("room_name", room.RoomName))
		}
	}

	logger.Info("Initial data seeded successfully")
}
