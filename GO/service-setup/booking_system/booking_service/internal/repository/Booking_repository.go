package repository

import (
	"context"
	"errors"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/errorcodes"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/pkg/database"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"booking_system/booking_service/internal/loggers"
)

type BookingRepositoryInterface interface {
	// Booking methods
	CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error
	GetAllBookings(ctx context.Context) ([]models.Booking, int64, *dtos.Error)
	GetBookingByUUID(ctx context.Context, bookingUUID string) (*models.Booking, *dtos.Error)
	PatchBooking(ctx context.Context, bookingUUID string, updateData map[string]interface{}) *dtos.Error
	CheckRoomAvailability(ctx context.Context, roomID int, bookingDate time.Time, newStartTime string, newEndTime string) (bool, *dtos.Error)

	// User methods
	GetUserByID(ctx context.Context, userID int) (*models.User, *dtos.Error)
	UserExists(ctx context.Context, userID int) (bool, *dtos.Error)

	// Room methods
	GetRoomByID(ctx context.Context, roomID int) (*models.Room, *dtos.Error)
	RoomExists(ctx context.Context, roomID int) (bool, *dtos.Error)

	// Seed methods
	SeedUsersAndRooms(ctx context.Context, users []models.User, rooms []models.Room) *dtos.Error
}

type BookingRepository struct {
	db     *database.Db
	logger *loggers.Logger
}

func NewBookingRepository(db *database.Db, logger *loggers.Logger) BookingRepositoryInterface {
	return &BookingRepository{
		db:     db,
		logger: logger,
	}
}

// ============ Booking Methods ============

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error {
	r.logger.Info("Creating new booking", zap.Int("user_id", booking.UserID), zap.Int("room_id", booking.RoomID))

	if err := r.db.Gorm.WithContext(ctx).Create(booking).Error; err != nil {
		r.logger.Error("Failed to create booking", zap.Error(err))
		return &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to create booking",
		}
	}

	r.logger.Info("Booking created successfully", zap.String("booking_uuid", booking.BookingUUID))
	return nil
}

func (r *BookingRepository) GetAllBookings(ctx context.Context) ([]models.Booking, int64, *dtos.Error) {
	r.logger.Info("Fetching all active bookings")

	var bookings []models.Booking
	var count int64

	// Count only active bookings
	if err := r.db.Gorm.WithContext(ctx).Where("is_active = ?", true).Model(&models.Booking{}).Count(&count).Error; err != nil {
		r.logger.Error("Failed to count bookings", zap.Error(err))
		return nil, 0, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to count bookings",
		}
	}

	// Fetch only active bookings
	if err := r.db.Gorm.WithContext(ctx).Where("is_active = ?", true).Find(&bookings).Error; err != nil {
		r.logger.Error("Failed to fetch bookings", zap.Error(err))
		return nil, 0, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to fetch bookings",
		}
	}

	r.logger.Info("Bookings fetched successfully", zap.Int64("count", count))
	return bookings, count, nil
}

func (r *BookingRepository) GetBookingByUUID(ctx context.Context, bookingUUID string) (*models.Booking, *dtos.Error) {
	r.logger.Info("Fetching booking by UUID", zap.String("booking_uuid", bookingUUID))

	var booking models.Booking
	if err := r.db.Gorm.WithContext(ctx).Where("booking_uuid = ?", bookingUUID).First(&booking).Error; err != nil {
		r.logger.Error("Failed to fetch booking", zap.String("booking_uuid", bookingUUID), zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "Booking not found",
			}
		}
		return nil, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to fetch booking",
		}
	}

	r.logger.Info("Booking fetched successfully", zap.String("booking_uuid", bookingUUID))
	return &booking, nil
}

func (r *BookingRepository) PatchBooking(ctx context.Context, bookingUUID string, updateData map[string]interface{}) *dtos.Error {
	r.logger.Info("Patching booking", zap.String("booking_uuid", bookingUUID))

	updateData["updated_at"] = time.Now().UTC()

	result := r.db.Gorm.WithContext(ctx).Model(&models.Booking{}).Where("booking_uuid = ?", bookingUUID).Updates(updateData)
	if result.Error != nil {
		r.logger.Error("Failed to patch booking", zap.String("booking_uuid", bookingUUID), zap.Error(result.Error))
		return &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to update booking",
		}
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("No booking found to update", zap.String("booking_uuid", bookingUUID))
		return &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
			Message: "Booking not found",
		}
	}

	r.logger.Info("Booking patched successfully", zap.String("booking_uuid", bookingUUID))
	return nil
}

// CheckRoomAvailability checks if room is available on given date
// Only active bookings block availability
func (r *BookingRepository) CheckRoomAvailability(ctx context.Context, roomID int, bookingDate time.Time, newStartTime string, newEndTime string) (bool, *dtos.Error) {
	r.logger.Info("Checking room availability", zap.Int("room_id", roomID), zap.String("booking_date", bookingDate.Format("2006-01-02")))

	var count int64
	// Check only for active bookings on same date
	if err := r.db.Gorm.WithContext(ctx).
		Where(`
		room_id = ?
		AND booking_date = ?
		AND is_active = ?
		AND ? < end_time
		AND ? > start_time
	`, roomID, bookingDate, true, newStartTime, newEndTime).
		Model(&models.Booking{}).
		Count(&count).Error; err != nil {
		r.logger.Error("Failed to check room availability", zap.Error(err))
		return false, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to check room availability",
		}
	}

	available := count == 0
	r.logger.Info("Room availability checked", zap.Int("room_id", roomID), zap.Bool("available", available))
	return available, nil
}

// ============ User Methods ============

func (r *BookingRepository) GetUserByID(ctx context.Context, userID int) (*models.User, *dtos.Error) {
	r.logger.Info("Fetching user by ID", zap.Int("user_id", userID))

	var user models.User
	if err := r.db.Gorm.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "User not found",
			}
		}
		r.logger.Error("Failed to fetch user", zap.Error(err))
		return nil, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to fetch user",
		}
	}

	return &user, nil
}

func (r *BookingRepository) UserExists(ctx context.Context, userID int) (bool, *dtos.Error) {
	var count int64
	if err := r.db.Gorm.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Model(&models.User{}).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check user existence", zap.Error(err))
		return false, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to check user existence",
		}
	}
	return count > 0, nil
}

// ============ Room Methods ============

func (r *BookingRepository) GetRoomByID(ctx context.Context, roomID int) (*models.Room, *dtos.Error) {
	r.logger.Info("Fetching room by ID", zap.Int("room_id", roomID))

	var room models.Room
	if err := r.db.Gorm.WithContext(ctx).Where("room_id = ? AND is_active = ?", roomID, true).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				Message: "Room not found",
			}
		}
		r.logger.Error("Failed to fetch room", zap.Error(err))
		return nil, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to fetch room",
		}
	}

	return &room, nil
}

func (r *BookingRepository) RoomExists(ctx context.Context, roomID int) (bool, *dtos.Error) {
	var count int64
	if err := r.db.Gorm.WithContext(ctx).Where("room_id = ? AND is_active = ?", roomID, true).Model(&models.Room{}).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check room existence", zap.Error(err))
		return false, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to check room existence",
		}
	}
	return count > 0, nil
}

// ============ Seed Methods ============

func (r *BookingRepository) SeedUsersAndRooms(ctx context.Context, users []models.User, rooms []models.Room) *dtos.Error {
	r.logger.Info("Seeding users and rooms")

	// Seed users
	for _, user := range users {
		if err := r.db.Gorm.WithContext(ctx).FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			r.logger.Error("Failed to seed user", zap.Error(err))
			return &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to seed users",
			}
		}
	}

	// Seed rooms
	for _, room := range rooms {
		if err := r.db.Gorm.WithContext(ctx).FirstOrCreate(&room, models.Room{RoomName: room.RoomName}).Error; err != nil {
			r.logger.Error("Failed to seed room", zap.Error(err))
			return &dtos.Error{
				Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				Message: "Failed to seed rooms",
			}
		}
	}

	r.logger.Info("Users and rooms seeded successfully")
	return nil
}
