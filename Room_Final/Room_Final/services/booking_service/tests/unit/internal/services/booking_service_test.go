package services_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/errorcodes"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/internal/services"
	repomock "booking_system/booking_service/tests/unit/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BookingServiceTestSuite struct {
	suite.Suite
	repoMock *repomock.BookingRepositoryMock
	service  services.BookingServiceInterface
	ctx      context.Context
}

func (suite *BookingServiceTestSuite) SetupTest() {
	suite.repoMock = new(repomock.BookingRepositoryMock)

	logger := loggers.NewLogger(loggers.LogConfig{
		Level:       "error",
		LogDir:      "./logs",
		FileName:    "service_test.log",
		ServiceName: "test-service",
	})

	suite.service = services.NewBookingService(suite.repoMock, logger)
	suite.ctx = context.WithValue(context.Background(), "user_id", "test-request-id-001")
}

func tomorrow() string {
	return time.Now().UTC().Add(48 * time.Hour).Format("2006-01-02")
}

// ---- CreateBooking ----

func (suite *BookingServiceTestSuite) TestCreateBooking_Success() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("GetRoomByID", suite.ctx, 1).Return(&models.Room{RoomID: 1, Capacity: 10}, nil)
	suite.repoMock.On("CheckRoomAvailability",
		suite.ctx, 1,
		mock.AnythingOfType("time.Time"),
		"09:00", "10:00", "").Return(true, nil)
	suite.repoMock.On("CreateBooking",
		suite.ctx,
		mock.AnythingOfType("*models.Booking")).Return(nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Success", response.Status)
	assert.Equal(suite.T(), "Booking created successfully", response.Message)
}

func (suite *BookingServiceTestSuite) TestCreateBooking_UserNotFound() {
	req := dtos.CreateBookingRequest{
		UserID:           999,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 999).Return(false, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "User not found")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_UserExistsDBError() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	dbErr := &dtos.Error{
		Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
		Message: "DB error",
	}
	suite.repoMock.On("UserExists", suite.ctx, 1).Return(false, dbErr)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Failed to validate user")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_RoomNotFound() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           999,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 999).Return(false, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Room not found")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_InvalidDateFormat() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      "31-12-2030",
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Invalid booking_date format")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_PastDate() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      "2020-01-01",
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "future")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_EndTimeBeforeStartTime() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "10:00",
		EndTime:          "09:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "End time must be greater than start time")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_ExceedsRoomCapacity() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 100,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("GetRoomByID", suite.ctx, 1).Return(&models.Room{RoomID: 1, Capacity: 10}, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "exceeds room capacity")
}

func (suite *BookingServiceTestSuite) TestCreateBooking_RoomNotAvailable() {
	req := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      tomorrow(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}

	suite.repoMock.On("UserExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("RoomExists", suite.ctx, 1).Return(true, nil)
	suite.repoMock.On("GetRoomByID", suite.ctx, 1).Return(&models.Room{RoomID: 1, Capacity: 10}, nil)
	suite.repoMock.On("CheckRoomAvailability",
		suite.ctx, 1,
		mock.AnythingOfType("time.Time"),
		"09:00", "10:00", "").Return(false, nil)

	response := suite.service.CreateBooking(suite.ctx, req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "already booked")
}

// ---- GetAllBookings ----

func (suite *BookingServiceTestSuite) TestGetAllBookings_Success() {
	now := time.Now()
	bookings := []models.Booking{
		{BookingUUID: "uuid-1", UserID: 1, RoomID: 1, BookingDate: now, StartTime: "09:00", EndTime: "10:00", IsActive: true},
		{BookingUUID: "uuid-2", UserID: 2, RoomID: 2, BookingDate: now, StartTime: "11:00", EndTime: "12:00", IsActive: true},
	}

	suite.repoMock.On("GetAllBookings", suite.ctx).Return(bookings, int64(2), nil)

	response := suite.service.GetAllBookings(suite.ctx)

	assert.Equal(suite.T(), "Success", response.Status)
	assert.NotNil(suite.T(), response.Data)

	listResp, ok := response.Data.(*dtos.BookingListResponse)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), int64(2), listResp.Total)
	assert.Len(suite.T(), listResp.Bookings, 2)
}

func (suite *BookingServiceTestSuite) TestGetAllBookings_DBError() {
	dbErr := &dtos.Error{
		Code:    errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
		Message: "DB failure",
	}

	suite.repoMock.On("GetAllBookings", suite.ctx).Return(nil, int64(0), dbErr)

	response := suite.service.GetAllBookings(suite.ctx)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Failed to fetch bookings")
}

func (suite *BookingServiceTestSuite) TestGetAllBookings_EmptyList() {
	suite.repoMock.On("GetAllBookings", suite.ctx).Return([]models.Booking{}, int64(0), nil)

	response := suite.service.GetAllBookings(suite.ctx)

	assert.Equal(suite.T(), "Success", response.Status)

	listResp, ok := response.Data.(*dtos.BookingListResponse)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), int64(0), listResp.Total)
	assert.Empty(suite.T(), listResp.Bookings)
}

// ---- GetBookingByUUID ----

func (suite *BookingServiceTestSuite) TestGetBookingByUUID_Success() {
	now := time.Now()
	booking := &models.Booking{
		BookingUUID: "test-uuid-001",
		UserID:      1,
		RoomID:      1,
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		IsActive:    true,
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(booking, nil)

	response := suite.service.GetBookingByUUID(suite.ctx, "test-uuid-001")

	assert.Equal(suite.T(), "Success", response.Status)
	assert.NotNil(suite.T(), response.Data)
}

func (suite *BookingServiceTestSuite) TestGetBookingByUUID_NotFound() {
	dbErr := &dtos.Error{
		Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
		Message: "Booking not found",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "non-existent").Return(nil, dbErr)

	response := suite.service.GetBookingByUUID(suite.ctx, "non-existent")

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Failed to fetch booking")
}

// ---- PatchBooking ----

func (suite *BookingServiceTestSuite) TestPatchBooking_Success_UpdateStartEndTime() {
	now := time.Now().Add(48 * time.Hour)
	newStart := "11:00"
	newEnd := "12:00"

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		UserID:      1,
		RoomID:      1,
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		StartTime: &newStart,
		EndTime:   &newEnd,
		UpdatedBy: "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)
	suite.repoMock.On("CheckRoomAvailability",
		suite.ctx, 1,
		mock.AnythingOfType("time.Time"),
		"11:00", "12:00", "test-uuid-001").Return(true, nil)
	suite.repoMock.On("PatchBooking",
		suite.ctx, "test-uuid-001",
		mock.AnythingOfType("map[string]interface {}")).Return(nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Success", response.Status)
}

func (suite *BookingServiceTestSuite) TestPatchBooking_BookingNotFound() {
	newStart := "11:00"
	req := dtos.PatchBookingRequest{
		StartTime: &newStart,
		UpdatedBy: "admin",
	}

	dbErr := &dtos.Error{
		Code:    errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
		Message: "Booking not found",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "non-existent").Return(nil, dbErr)

	response := suite.service.PatchBooking(suite.ctx, "non-existent", req)

	assert.Equal(suite.T(), "Error", response.Status)
}

func (suite *BookingServiceTestSuite) TestPatchBooking_CancelBooking_Success() {
	now := time.Now().Add(48 * time.Hour)
	isActive := false

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		IsActive:  &isActive,
		UpdatedBy: "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)
	suite.repoMock.On("PatchBooking",
		suite.ctx, "test-uuid-001",
		mock.AnythingOfType("map[string]interface {}")).Return(nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Success", response.Status)
}

func (suite *BookingServiceTestSuite) TestPatchBooking_AlreadyActive_Error() {
	now := time.Now().Add(48 * time.Hour)
	isActive := true

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		IsActive:  &isActive,
		UpdatedBy: "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "already active")
}

func (suite *BookingServiceTestSuite) TestPatchBooking_AlreadyCancelled_Error() {
	now := time.Now().Add(48 * time.Hour)
	isActive := false

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    false,
	}

	req := dtos.PatchBookingRequest{
		IsActive:  &isActive,
		UpdatedBy: "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "already cancelled")
}

func (suite *BookingServiceTestSuite) TestPatchBooking_InvalidDateFormat() {
	now := time.Now().Add(48 * time.Hour)
	badDate := "31-12-2030"

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		BookingDate: &badDate,
		UpdatedBy:   "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "Invalid booking_date format")
}

func (suite *BookingServiceTestSuite) TestPatchBooking_EndTimeBeforeStart() {
	now := time.Now().Add(48 * time.Hour)
	newEnd := "08:00"

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		EndTime:   &newEnd,
		UpdatedBy: "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "End time must be greater than start time")
}

func (suite *BookingServiceTestSuite) TestPatchBooking_ExceedsCapacity() {
	now := time.Now().Add(48 * time.Hour)
	newCount := 100

	existing := &models.Booking{
		BookingUUID: "test-uuid-001",
		BookingDate: now,
		StartTime:   "09:00",
		EndTime:     "10:00",
		RoomID:      1,
		IsActive:    true,
	}

	req := dtos.PatchBookingRequest{
		NumberOfEmployee: &newCount,
		UpdatedBy:        "admin",
	}

	suite.repoMock.On("GetBookingByUUID", suite.ctx, "test-uuid-001").Return(existing, nil)
	suite.repoMock.On("GetRoomByID", suite.ctx, 1).Return(&models.Room{RoomID: 1, Capacity: 10}, nil)

	response := suite.service.PatchBooking(suite.ctx, "test-uuid-001", req)

	assert.Equal(suite.T(), "Error", response.Status)
	assert.Contains(suite.T(), response.Message, "exceeds room capacity")
}

// unused import guard
var _ = http.StatusOK

func TestBookingServiceSuite(t *testing.T) {
	suite.Run(t, new(BookingServiceTestSuite))
}
