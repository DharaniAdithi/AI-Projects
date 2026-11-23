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
			Errors:    []dtos.Error{{Field: "booking_date", Message: "Wrong Format", Code: "TRAINING_VAL_001"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Future date validation
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if !bookingDate.After(today) {
		s.logger.Warn("Booking date is not in the future", zap.String("booking_date", req.BookingDate))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Booking date must be in the future",
			Errors:    []dtos.Error{{Field: "booking_date", Message: "Booking date must be in the future", Code: "TRAINING_VAL_001"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Time validation
	if req.EndTime <= req.StartTime {
		s.logger.Warn("End time is not after start time", zap.String("start_time", req.StartTime), zap.String("end_time", req.EndTime))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "End time must be greater than start time",
			Errors:    []dtos.Error{{Field: "end_time", Message: "End time must be greater than start time", Code: "TRAINING_VAL_001"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Get room for capacity check
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

	// Room capacity validation
	if req.NumberOfEmployee > room.Capacity {
		s.logger.Warn("Number of employees exceeds room capacity", zap.Int("employees", req.NumberOfEmployee), zap.Int("capacity", room.Capacity))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Number of employees exceeds room capacity",
			Errors:    []dtos.Error{{Field: "number_of_employee", Message: "Number of employees exceeds room capacity", Code: "TRAINING_VAL_001"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Check room availability
	available, availErr := s.bookingRepo.CheckRoomAvailability(ctx, req.RoomID, bookingDate, req.StartTime, req.EndTime, "")
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
			Message:   "Room is already booked for the selected time",
			Errors:    []dtos.Error{{Field: "room_id", Message: "Room is already booked for the selected time", Code: "TRAINING_VAL_003"}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	booking := &models.Booking{
		BookingUUID:   ctx.Value("user_id").(string),
		UserID:        req.UserID,
		RoomID:        req.RoomID,
		BookingDate:   bookingDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
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

	// Fetch existing booking
	existing, repoErr := s.bookingRepo.GetBookingByUUID(ctx, bookingUUID)
	if repoErr != nil {
		s.logger.Error("Booking not found for patch", zap.String("booking_uuid", bookingUUID))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Failed to fetch booking",
			Errors:    []dtos.Error{{Message: repoErr.Message, Code: repoErr.Code}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Guard against redundant is_active toggles
	if req.IsActive != nil {
		if *req.IsActive && existing.IsActive {
			s.logger.Warn("Booking is already active", zap.String("booking_uuid", bookingUUID))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Booking is already active",
				Errors:    []dtos.Error{{Field: "is_active", Message: "Booking is already active", Code: "TRAINING_VAL_003"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		if !*req.IsActive && !existing.IsActive {
			s.logger.Warn("Booking is already cancelled", zap.String("booking_uuid", bookingUUID))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Booking is already cancelled",
				Errors:    []dtos.Error{{Field: "is_active", Message: "Booking is already cancelled", Code: "TRAINING_VAL_003"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	// Compute effective values by merging request over existing
	effectiveDate := existing.BookingDate
	effectiveStart := existing.StartTime
	effectiveEnd := existing.EndTime
	effectiveRoomID := existing.RoomID

	if req.BookingDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.BookingDate)
		if err != nil {
			s.logger.Warn("Invalid booking date format", zap.String("booking_date", *req.BookingDate))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Invalid booking_date format. Expected YYYY-MM-DD",
				Errors:    []dtos.Error{{Field: "booking_date", Message: "Wrong Format", Code: "TRAINING_VAL_001"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		today := time.Now().UTC().Truncate(24 * time.Hour)
		if !parsedDate.After(today) {
			s.logger.Warn("Booking date is not in the future", zap.String("booking_date", *req.BookingDate))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Booking date must be in the future",
				Errors:    []dtos.Error{{Field: "booking_date", Message: "Booking date must be in the future", Code: "TRAINING_VAL_001"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		effectiveDate = parsedDate
	}

	if req.StartTime != nil {
		effectiveStart = *req.StartTime
	}
	if req.EndTime != nil {
		effectiveEnd = *req.EndTime
	}
	if req.RoomID != nil {
		effectiveRoomID = *req.RoomID
	}

	// Time validation when start or end changes
	if req.StartTime != nil || req.EndTime != nil {
		if effectiveEnd <= effectiveStart {
			s.logger.Warn("End time is not after start time", zap.String("start_time", effectiveStart), zap.String("end_time", effectiveEnd))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "End time must be greater than start time",
				Errors:    []dtos.Error{{Field: "end_time", Message: "End time must be greater than start time", Code: "TRAINING_VAL_001"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	// Room capacity validation when employee count is provided
	if req.NumberOfEmployee != nil {
		room, roomErr := s.bookingRepo.GetRoomByID(ctx, effectiveRoomID)
		if roomErr != nil {
			s.logger.Error("Failed to fetch room for capacity check", zap.String("error", roomErr.Message))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Failed to fetch room details",
				Errors:    []dtos.Error{{Message: roomErr.Message, Code: roomErr.Code}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		if *req.NumberOfEmployee > room.Capacity {
			s.logger.Warn("Number of employees exceeds room capacity", zap.Int("employees", *req.NumberOfEmployee), zap.Int("capacity", room.Capacity))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Number of employees exceeds room capacity",
				Errors:    []dtos.Error{{Field: "number_of_employee", Message: "Number of employees exceeds room capacity", Code: "TRAINING_VAL_001"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	// Availability check: needed when scheduling fields change, or when recovering
	isRecovering := req.IsActive != nil && *req.IsActive && !existing.IsActive
	schedulingChanged := req.BookingDate != nil || req.StartTime != nil || req.EndTime != nil || req.RoomID != nil

	if schedulingChanged || isRecovering {
		available, availErr := s.bookingRepo.CheckRoomAvailability(ctx, effectiveRoomID, effectiveDate, effectiveStart, effectiveEnd, bookingUUID)
		if availErr != nil {
			s.logger.Error("Failed to check room availability", zap.String("error", availErr.Message))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Failed to check room availability",
				Errors:    []dtos.Error{{Message: availErr.Message, Code: availErr.Code}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
		if !available {
			s.logger.Warn("Room not available for updated time", zap.Int("room_id", effectiveRoomID))
			return &dtos.APIResponse{
				Status:    "Error",
				Message:   "Room is already booked for the selected time",
				Errors:    []dtos.Error{{Field: "room_id", Message: "Room is already booked for the selected time", Code: "TRAINING_VAL_003"}},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	// Build update map
	updateData := map[string]interface{}{
		"updated_by": req.UpdatedBy,
	}

	if req.BookingDate != nil {
		updateData["booking_date"] = effectiveDate
	}
	if req.StartTime != nil {
		updateData["start_time"] = effectiveStart
	}
	if req.EndTime != nil {
		updateData["end_time"] = effectiveEnd
	}
	if req.RoomID != nil {
		updateData["room_id"] = effectiveRoomID
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
		if *req.IsActive {
			updateData["booking_status"] = "active"
		} else {
			updateData["booking_status"] = "cancelled"
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

func mapBookingToResponse(b *models.Booking) dtos.BookingResponse {
	status := "Active"
	if !b.IsActive {
		status = "Cancelled"
	}
	return dtos.BookingResponse{
		BookingUUID:   b.BookingUUID,
		UserID:        b.UserID,
		RoomID:        b.RoomID,
		BookingDate:   b.BookingDate.Format("2006-01-02"),
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		BookingStatus: status,
		IsActive:      b.IsActive,
		CreatedBy:     b.CreatedBy,
		UpdatedBy:     b.UpdatedBy,
		CreatedAt:     b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
