package repository_test


import (
	"context"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/models"

	"github.com/stretchr/testify/mock"
)

// BookingRepositoryMock mocks BookingRepositoryInterface
type BookingRepositoryMock struct {
	mock.Mock
}

func (m *BookingRepositoryMock) CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error {
	args := m.Called(ctx, booking)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dtos.Error)
}

func (m *BookingRepositoryMock) GetAllBookings(ctx context.Context) ([]models.Booking, int64, *dtos.Error) {
	args := m.Called(ctx)
	var bookings []models.Booking
	if args.Get(0) != nil {
		bookings = args.Get(0).([]models.Booking)
	}
	count := args.Get(1).(int64)
	var err *dtos.Error
	if args.Get(2) != nil {
		err = args.Get(2).(*dtos.Error)
	}
	return bookings, count, err
}

func (m *BookingRepositoryMock) GetBookingByUUID(ctx context.Context, bookingUUID string) (*models.Booking, *dtos.Error) {
	args := m.Called(ctx, bookingUUID)
	var booking *models.Booking
	if args.Get(0) != nil {
		booking = args.Get(0).(*models.Booking)
	}
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return booking, err
}

func (m *BookingRepositoryMock) PatchBooking(ctx context.Context, bookingUUID string, updateData map[string]interface{}) *dtos.Error {
	args := m.Called(ctx, bookingUUID, updateData)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dtos.Error)
}

func (m *BookingRepositoryMock) CheckRoomAvailability(ctx context.Context, roomID int, bookingDate time.Time, newStartTime string, newEndTime string, excludeBookingUUID string) (bool, *dtos.Error) {
	args := m.Called(ctx, roomID, bookingDate, newStartTime, newEndTime, excludeBookingUUID)
	available := args.Bool(0)
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return available, err
}

func (m *BookingRepositoryMock) GetUserByID(ctx context.Context, userID int) (*models.User, *dtos.Error) {
	args := m.Called(ctx, userID)
	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return user, err
}

func (m *BookingRepositoryMock) UserExists(ctx context.Context, userID int) (bool, *dtos.Error) {
	args := m.Called(ctx, userID)
	exists := args.Bool(0)
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return exists, err
}

func (m *BookingRepositoryMock) GetRoomByID(ctx context.Context, roomID int) (*models.Room, *dtos.Error) {
	args := m.Called(ctx, roomID)
	var room *models.Room
	if args.Get(0) != nil {
		room = args.Get(0).(*models.Room)
	}
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return room, err
}

func (m *BookingRepositoryMock) RoomExists(ctx context.Context, roomID int) (bool, *dtos.Error) {
	args := m.Called(ctx, roomID)
	exists := args.Bool(0)
	var err *dtos.Error
	if args.Get(1) != nil {
		err = args.Get(1).(*dtos.Error)
	}
	return exists, err
}


