package services

import (
	"context"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/internal/repository"

	"go.uber.org/zap"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, req dtos.CreateBookingRequest) *dtos.APIResponse
	GetAllBookings(ctx context.Context) *dtos.APIResponse
	GetBookingByUUID(ctx context.Context, bookingUUID string) *dtos.APIResponse
	PatchBooking(ctx context.Context, bookingUUID string, req dtos.PatchBookingRequest) *dtos.APIResponse
	CancelBooking(ctx context.Context, bookingUUID string, req dtos.CancelBookingRequest) *dtos.APIResponse
	RecoverBooking(ctx context.Context, bookingUUID string, req dtos.RecoverBookingRequest) *dtos.APIResponse
	SeedUsersAndRooms(ctx context.Context, users []models.User, rooms []models.Room) *dtos.APIResponse
}

type BookingService struct {
	bookingRepo repository.BookingRepositoryInterface
	logger      *loggers.Logger
}

func NewBookingService(bookingRepo repository.BookingRepositoryInterface, logger *loggers.Logger) BookingServiceInterface {
	return &BookingService{
		bookingRepo: bookingRepo,
		logger:      logger,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req dtos.CreateBookingRequest) *dtos.APIResponse {
	s.logger.Info("Creating new booking", zap.Int("user_id", req.UserID), zap.Int("room_id", req.RoomID))

	// Validate user exists
	userExists, err := s.bookingRepo.UserExists(ctx, req.UserID)
	if err != nil {
		s.logger.Error("Failed to validate user", zap.String("error", err.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to validate user",
			Errors:    []dtos.Error{{Message: err.Message, Code: err.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}
	if !userExists {
		s.logger.Warn("User not found", zap.Int("user_id", req.UserID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "User not found",
			Errors:    []dtos.Error{{Field: "user_id", Message: "User does not exist", Code: "TRAINING_REC_404"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Validate room exists
	roomExists, err := s.bookingRepo.RoomExists(ctx, req.RoomID)
	if err != nil {
		s.logger.Error("Failed to validate room", zap.String("error", err.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to validate room",
			Errors:    []dtos.Error{{Message: err.Message, Code: err.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}
	if !roomExists {
		s.logger.Warn("Room not found", zap.Int("room_id", req.RoomID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Room not found",
			Errors:    []dtos.Error{{Field: "room_id", Message: "Room does not exist", Code: "TRAINING_REC_404"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Parse booking date
	bookingDate, parseErr := time.Parse("2006-01-02", req.BookingDate)
	if parseErr != nil {
		s.logger.Warn("Invalid booking date format", zap.String("booking_date", req.BookingDate))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Invalid booking_date format. Expected YYYY-MM-DD",
			Errors:    []dtos.Error{{Message: "Wrong Format", Code: "TRAINING_VAL_001"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Check room availability
	available, availErr := s.bookingRepo.CheckRoomAvailability(ctx, req.RoomID, bookingDate, req.StartTime, req.EndTime)
	if availErr != nil {
		s.logger.Error("Failed to check availability", zap.String("error", availErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to check room availability",
			Errors:    []dtos.Error{{Message: availErr.Message, Code: availErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}
	if !available {
		s.logger.Warn("Room not available", zap.Int("room_id", req.RoomID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Room is not available on the selected date",
			Errors:    []dtos.Error{{Field: "room_id", Message: "Room already booked on this date", Code: "TRAINING_VAL_003"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Get room pricing
	room, roomErr := s.bookingRepo.GetRoomByID(ctx, req.RoomID)
	if roomErr != nil {
		s.logger.Error("Failed to fetch room", zap.String("error", roomErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to fetch room details",
			Errors:    []dtos.Error{{Message: roomErr.Message, Code: roomErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Create booking
	booking := &models.Booking{
		BookingUUID:   ctx.Value("user_id").(string),
		UserID:        req.UserID,
		RoomID:        req.RoomID,
		BookingDate:   bookingDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		PricePerDay:   room.PricePerDay, // Capture price at booking time
		BookingStatus: "active",
		IsActive:      true,
		CreatedBy:     req.CreatedBy,
		UpdatedBy:     req.CreatedBy,
	}

	if repoErr := s.bookingRepo.CreateBooking(ctx, booking); repoErr != nil {
		s.logger.Error("Failed to create booking", zap.String("error", repoErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to create booking",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	response := mapBookingToResponse(booking)

	s.logger.Info("Booking created successfully", zap.String("booking_uuid", booking.BookingUUID))
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Booking created successfully",
		Data:      response,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *BookingService) GetAllBookings(ctx context.Context) *dtos.APIResponse {
	s.logger.Info("Fetching all bookings")

	bookings, count, repoErr := s.bookingRepo.GetAllBookings(ctx)
	if repoErr != nil {
		s.logger.Error("Failed to fetch all bookings", zap.String("error", repoErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to fetch bookings",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	bookingResponses := make([]dtos.BookingResponse, len(bookings))
	for i, b := range bookings {
		bookingResponses[i] = mapBookingToResponse(&b)
	}

	s.logger.Info("All bookings fetched successfully", zap.Int64("count", count))
	return &dtos.APIResponse{
		Status:  "Success",
		Message: "Bookings fetched successfully",
		Data: &dtos.BookingListResponse{
			Bookings: bookingResponses,
			Total:    count,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *BookingService) GetBookingByUUID(ctx context.Context, bookingUUID string) *dtos.APIResponse {
	s.logger.Info("Fetching booking by UUID", zap.String("booking_uuid", bookingUUID))

	booking, repoErr := s.bookingRepo.GetBookingByUUID(ctx, bookingUUID)
	if repoErr != nil {
		s.logger.Error("Failed to fetch booking", zap.String("booking_uuid", bookingUUID), zap.String("error", repoErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to fetch booking",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	response := mapBookingToResponse(booking)

	s.logger.Info("Booking fetched successfully", zap.String("booking_uuid", bookingUUID))
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Booking fetched successfully",
		Data:      response,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *BookingService) PatchBooking(ctx context.Context, bookingUUID string, req dtos.PatchBookingRequest) *dtos.APIResponse {
	s.logger.Info("Patching booking", zap.String("booking_uuid", bookingUUID))

	// Verify booking exists
	_, repoErr := s.bookingRepo.GetBookingByUUID(ctx, bookingUUID)
	if repoErr != nil {
		s.logger.Error("Booking not found for patch", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to fetch booking",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	updateData := map[string]interface{}{
		"updated_by": req.UpdatedBy,
	}

	if req.BookingDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.BookingDate)
		if err != nil {
			s.logger.Warn("Invalid booking date format", zap.String("booking_date", *req.BookingDate))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Invalid booking_date format. Expected YYYY-MM-DD",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		updateData["booking_date"] = parsedDate
	}
	if req.StartTime != nil {
		updateData["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		updateData["end_time"] = *req.EndTime
	}
	if req.BookingStatus != nil {
		updateData["booking_status"] = *req.BookingStatus
		// Soft delete: if status is cancelled, mark is_active false
		if *req.BookingStatus == "cancelled" {
			updateData["is_active"] = false
		}
	}

	if patchErr := s.bookingRepo.PatchBooking(ctx, bookingUUID, updateData); patchErr != nil {
		s.logger.Error("Failed to patch booking", zap.String("booking_uuid", bookingUUID), zap.String("error", patchErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to update booking",
			Errors:    []dtos.Error{{Message: patchErr.Message, Code: patchErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	s.logger.Info("Booking patched successfully", zap.String("booking_uuid", bookingUUID))
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Booking updated successfully",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CancelBooking cancels a booking (soft delete)
func (s *BookingService) CancelBooking(ctx context.Context, bookingUUID string, req dtos.CancelBookingRequest) *dtos.APIResponse {
	s.logger.Info("Cancelling booking", zap.String("booking_uuid", bookingUUID))

	// Verify booking exists and is active
	booking, repoErr := s.bookingRepo.GetBookingByUUID(ctx, bookingUUID)
	if repoErr != nil {
		s.logger.Error("Booking not found", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Booking not found",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	if !booking.IsActive {
		s.logger.Warn("Booking already cancelled", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Booking is already cancelled",
			Errors:    []dtos.Error{{Field: "booking_uuid", Message: "Booking cannot be cancelled", Code: "TRAINING_VAL_003"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	updateData := map[string]interface{}{
		"booking_status": "cancelled",
		"is_active":      false,
		"updated_by":     req.UpdatedBy,
	}

	if patchErr := s.bookingRepo.PatchBooking(ctx, bookingUUID, updateData); patchErr != nil {
		s.logger.Error("Failed to cancel booking", zap.String("booking_uuid", bookingUUID), zap.String("error", patchErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to cancel booking",
			Errors:    []dtos.Error{{Message: patchErr.Message, Code: patchErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	s.logger.Info("Booking cancelled successfully", zap.String("booking_uuid", bookingUUID))
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Booking cancelled successfully. Room is now available for booking.",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// CancelBooking cancels a booking (soft delete)
func (s *BookingService) RecoverBooking(ctx context.Context, bookingUUID string, req dtos.RecoverBookingRequest) *dtos.APIResponse {
	s.logger.Info("Recovering booking", zap.String("booking_uuid", bookingUUID))

	// Verify booking exists and is active
	booking, repoErr := s.bookingRepo.GetBookingByUUID(ctx, bookingUUID)
	if repoErr != nil {
		s.logger.Error("Booking not found", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Booking not found",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	if booking.IsActive {
		s.logger.Warn("Booking already active", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Booking is already active",
			Errors:    []dtos.Error{{Field: "booking_uuid", Message: "Booking cannot be cancelled", Code: "TRAINING_VAL_003"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	roomIsBookedBySomeOne, err := s.bookingRepo.CheckRoomAvailability(ctx, booking.RoomID, booking.BookingDate, booking.StartTime, booking.EndTime)

	if roomIsBookedBySomeOne {
		s.logger.Error("Failed to recover booking", zap.String("booking_uuid", bookingUUID), zap.String("error", err.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to recover booking",
			Errors:    []dtos.Error{{Message: err.Message, Code: err.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	updateData := map[string]interface{}{
		"booking_status": "active",
		"is_active":      true,
		"updated_by":     req.UpdatedBy,
	}

	if patchErr := s.bookingRepo.PatchBooking(ctx, bookingUUID, updateData); patchErr != nil {
		s.logger.Error("Failed to recover booking", zap.String("booking_uuid", bookingUUID), zap.String("error", patchErr.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to recover booking",
			Errors:    []dtos.Error{{Message: patchErr.Message, Code: patchErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	s.logger.Info("Booking recovered successfully", zap.String("booking_uuid", bookingUUID))
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Booking recovered successfully.",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// SeedUsersAndRooms seeds initial users and rooms
func (s *BookingService) SeedUsersAndRooms(ctx context.Context, users []models.User, rooms []models.Room) *dtos.APIResponse {
	s.logger.Info("Seeding users and rooms", zap.Int("users_count", len(users)), zap.Int("rooms_count", len(rooms)))

	if err := s.bookingRepo.SeedUsersAndRooms(ctx, users, rooms); err != nil {
		s.logger.Error("Failed to seed data", zap.String("error", err.Message))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to seed users and rooms",
			Errors:    []dtos.Error{{Message: err.Message, Code: err.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	s.logger.Info("Data seeded successfully")
	return &dtos.APIResponse{
		Status:    "Success",
		Message:   "Users and rooms seeded successfully",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func mapBookingToResponse(b *models.Booking) dtos.BookingResponse {
	return dtos.BookingResponse{
		BookingUUID:   b.BookingUUID,
		UserID:        b.UserID,
		RoomID:        b.RoomID,
		BookingDate:   b.BookingDate.Format("2006-01-02"),
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		PricePerDay:   b.PricePerDay,
		BookingStatus: b.BookingStatus,
		IsActive:      b.IsActive,
		CreatedBy:     b.CreatedBy,
		UpdatedBy:     b.UpdatedBy,
		CreatedAt:     b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
