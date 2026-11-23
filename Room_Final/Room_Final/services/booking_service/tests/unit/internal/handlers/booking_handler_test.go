package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/errorcodes"
	"booking_system/booking_service/internal/handlers"
	"booking_system/booking_service/internal/loggers"
	svcmock "booking_system/booking_service/tests/unit/internal/services/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BookingHandlerTestSuite struct {
	suite.Suite
	serviceMock *svcmock.BookingServiceMock
	router      *mux.Router
}

func (suite *BookingHandlerTestSuite) SetupTest() {
	suite.serviceMock = new(svcmock.BookingServiceMock)

	logger := loggers.NewLogger(loggers.LogConfig{
		Level:       "info",
		LogDir:      "./logs",
		FileName:    "test.log",
		ServiceName: "test-service",
	})

	handler := handlers.NewBookingHandler(suite.serviceMock, logger)
	suite.router = mux.NewRouter()
	suite.router.HandleFunc("/api/v1/booking", handler.CreateBooking).Methods("POST")
	suite.router.HandleFunc("/api/v1/bookings", handler.GetAllBookings).Methods("GET")
	suite.router.HandleFunc("/api/v1/bookings/{uuid}", handler.GetBookingByUUID).Methods("GET")
	suite.router.HandleFunc("/api/v1/bookings/{uuid}", handler.PatchBooking).Methods("PATCH")
}

// ---- CreateBooking Handler ----

func (suite *BookingHandlerTestSuite) TestCreateBooking_Success() {
	payload := dtos.CreateBookingRequest{
		UserID:           1,
		RoomID:           1,
		BookingDate:      time.Now().Add(48 * time.Hour).Format("2006-01-02"),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		CreatedBy:        "admin",
	}
	body, _ := json.Marshal(payload)

	suite.serviceMock.On("CreateBooking", mock.Anything, mock.AnythingOfType("dtos.CreateBookingRequest")).
		Return(&dtos.APIResponse{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: "Booking created successfully",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("POST", "/api/v1/booking", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-request-id")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var resp dtos.APIResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(suite.T(), "Success", resp.Status)
}

func (suite *BookingHandlerTestSuite) TestCreateBooking_InvalidJSON() {
	body := []byte(`{"user_id": not-a-number}`) // invalid type for int field

	req, _ := http.NewRequest("POST", "/api/v1/booking", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *BookingHandlerTestSuite) TestCreateBooking_ValidationFailed() {
	payload := map[string]interface{}{
		"user_id": 0,  // fails min=1 validation
		"room_id": 1,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/v1/booking", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

// ---- GetAllBookings Handler ----

func (suite *BookingHandlerTestSuite) TestGetAllBookings_Success() {
	suite.serviceMock.On("GetAllBookings", mock.Anything).
		Return(&dtos.APIResponse{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: "Bookings fetched successfully",
			Data: &dtos.BookingListResponse{
				Bookings: []dtos.BookingResponse{},
				Total:    0,
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("GET", "/api/v1/bookings", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var resp dtos.APIResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(suite.T(), "Success", resp.Status)
}

func (suite *BookingHandlerTestSuite) TestGetAllBookings_ServiceError() {
	suite.serviceMock.On("GetAllBookings", mock.Anything).
		Return(&dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch bookings",
			Errors: []dtos.Error{
				{Code: errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode], Message: "DB error"},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("GET", "/api/v1/bookings", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
}

// ---- GetBookingByUUID Handler ----

func (suite *BookingHandlerTestSuite) TestGetBookingByUUID_Success() {
	suite.serviceMock.On("GetBookingByUUID", mock.Anything, "test-uuid-001").
		Return(&dtos.APIResponse{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: "Booking fetched successfully",
			Data: dtos.BookingResponse{
				BookingUUID: "test-uuid-001",
				UserID:      1,
				RoomID:      1,
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("GET", "/api/v1/bookings/test-uuid-001", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func (suite *BookingHandlerTestSuite) TestGetBookingByUUID_NotFound() {
	suite.serviceMock.On("GetBookingByUUID", mock.Anything, "non-existent").
		Return(&dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusNotFound,
			Message: "Failed to fetch booking",
			Errors: []dtos.Error{
				{Code: errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode], Message: "Booking not found"},
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("GET", "/api/v1/bookings/non-existent", nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusNotFound, rr.Code)
}

// ---- PatchBooking Handler ----

func (suite *BookingHandlerTestSuite) TestPatchBooking_Success() {
	newStart := "11:00"
	payload := dtos.PatchBookingRequest{
		StartTime: &newStart,
		UpdatedBy: "admin",
	}
	body, _ := json.Marshal(payload)

	suite.serviceMock.On("PatchBooking", mock.Anything, "test-uuid-001",
		mock.AnythingOfType("dtos.PatchBookingRequest")).
		Return(&dtos.APIResponse{
			Status:    "Success",
			Code:      http.StatusOK,
			Message:   "Booking updated successfully",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})

	req, _ := http.NewRequest("PATCH", "/api/v1/bookings/test-uuid-001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
}

func (suite *BookingHandlerTestSuite) TestPatchBooking_InvalidJSON() {
	body := []byte(`{invalid json}`)

	req, _ := http.NewRequest("PATCH", "/api/v1/bookings/test-uuid-001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *BookingHandlerTestSuite) TestPatchBooking_ValidationFailed_MissingUpdatedBy() {
	payload := map[string]interface{}{
		"start_time": "11:00",
		// updated_by missing - required field
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "/api/v1/bookings/test-uuid-001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func TestBookingHandlerSuite(t *testing.T) {
	suite.Run(t, new(BookingHandlerTestSuite))
}


