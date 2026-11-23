package services

import (
	"context"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/internal/repository"

	"github.com/google/uuid"
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
	s.logger.Info("Creating new booking", zap.Int("guest_id", req.GuestID), zap.Int("room_id", req.RoomID))

	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		s.logger.Warn("Invalid booking date format", zap.String("booking_date", req.BookingDate))
		return &dtos.APIResponse{
			Status:    "Error",
			Message:   "Invalid booking_date format. Expected YYYY-MM-DD",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	booking := &models.Booking{
		BookingUUID:   uuid.NewString(),
		GuestID:       req.GuestID,
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
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
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
	return dtos.BookingResponse{
		BookingUUID:   b.BookingUUID,
		GuestID:       b.GuestID,
		RoomID:        b.RoomID,
		BookingDate:   b.BookingDate.Format("2006-01-02"),
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		BookingStatus: b.BookingStatus,
		IsActive:      b.IsActive,
		CreatedBy:     b.CreatedBy,
		UpdatedBy:     b.UpdatedBy,
		CreatedAt:     b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
