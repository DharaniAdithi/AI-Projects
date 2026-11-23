package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/errorcodes"
	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/services"
	"booking_system/booking_service/internal/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type BookingHandler struct {
	bookingService services.BookingServiceInterface
	logger         *loggers.Logger
}

type BookingHandlerInterface interface {
	CreateBooking(w http.ResponseWriter, r *http.Request)
	GetAllBookings(w http.ResponseWriter, r *http.Request)
	GetBookingByUUID(w http.ResponseWriter, r *http.Request)
	PatchBooking(w http.ResponseWriter, r *http.Request)
	CancelBooking(w http.ResponseWriter, r *http.Request)
	RecoverBooking(w http.ResponseWriter, r *http.Request)
}

func NewBookingHandler(bookingService services.BookingServiceInterface, logger *loggers.Logger) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		logger:         logger,
	}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to create booking")
	requestId := r.Header.Get("X-Request-ID")

	ctx := context.WithValue(r.Context(), "user_id", requestId)

	r = r.WithContext(ctx)

	var req dtos.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		h.logger.Info("Validation failed for create booking request")
		apiResponse := &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			Message:   "Validation failed for the request",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	apiResponse := h.bookingService.CreateBooking(r.Context(), req)
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}

func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to get all bookings")

	apiResponse := h.bookingService.GetAllBookings(r.Context())
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}

func (h *BookingHandler) GetBookingByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingUUID := vars["uuid"]

	h.logger.Info("Received request to get booking", zap.String("booking_uuid", bookingUUID))

	if bookingUUID == "" {
		h.logger.Info("Missing UUID in request")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Validation failed for the request",
			Errors: []dtos.Error{{
				Field:   "uuid",
				Message: "Booking UUID is required",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	apiResponse := h.bookingService.GetBookingByUUID(r.Context(), bookingUUID)
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}

func (h *BookingHandler) PatchBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingUUID := vars["uuid"]

	h.logger.Info("Received request to patch booking", zap.String("booking_uuid", bookingUUID))

	if bookingUUID == "" {
		h.logger.Info("Missing UUID in request")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Validation failed for the request",
			Errors: []dtos.Error{{
				Field:   "uuid",
				Message: "Booking UUID is required",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	var req dtos.PatchBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		h.logger.Info("Validation failed for patch booking request")
		apiResponse := &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			Message:   "Validation failed for the request",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	apiResponse := h.bookingService.PatchBooking(r.Context(), bookingUUID, req)
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}

// CancelBooking cancels a booking with soft delete
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingUUID := vars["uuid"]

	h.logger.Info("Received request to cancel booking", zap.String("booking_uuid", bookingUUID))

	if bookingUUID == "" {
		h.logger.Info("Missing UUID in request")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Validation failed for the request",
			Errors: []dtos.Error{{
				Field:   "uuid",
				Message: "Booking UUID is required",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	var req dtos.CancelBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		h.logger.Info("Validation failed for cancel booking request")
		apiResponse := &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			Message:   "Validation failed for the request",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	apiResponse := h.bookingService.CancelBooking(r.Context(), bookingUUID, req)
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}



// CancelBooking cancels a booking with soft delete
func (h *BookingHandler) RecoverBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingUUID := vars["uuid"]

	h.logger.Info("Received request to cancel booking", zap.String("booking_uuid", bookingUUID))

	if bookingUUID == "" {
		h.logger.Info("Missing UUID in request")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Validation failed for the request",
			Errors: []dtos.Error{{
				Field:   "uuid",
				Message: "Booking UUID is required",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	var req dtos.RecoverBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode],
			}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		h.logger.Info("Validation failed for cancel booking request")
		apiResponse := &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			Message:   "Validation failed for the request",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	apiResponse := h.bookingService.RecoverBooking(r.Context(), bookingUUID, req)
	response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, apiResponse)
}
