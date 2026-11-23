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
	CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error
	GetAllBookings(ctx context.Context) ([]models.Booking, int64, *dtos.Error)
	GetBookingByUUID(ctx context.Context, bookingUUID string) (*models.Booking, *dtos.Error)
	PatchBooking(ctx context.Context, bookingUUID string, updateData map[string]interface{}) *dtos.Error
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

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error {
	r.logger.Info("Creating new booking", zap.Int("guest_id", booking.GuestID), zap.Int("room_id", booking.RoomID))

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
	r.logger.Info("Fetching all bookings")

	var bookings []models.Booking
	var count int64

	if err := r.db.Gorm.WithContext(ctx).Model(&models.Booking{}).Count(&count).Error; err != nil {
		r.logger.Error("Failed to count bookings", zap.Error(err))
		return nil, 0, &dtos.Error{
			Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
			Message: "Failed to count bookings",
		}
	}

	if err := r.db.Gorm.WithContext(ctx).Find(&bookings).Error; err != nil {
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
