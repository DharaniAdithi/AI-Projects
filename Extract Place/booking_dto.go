package dtos

type CreateBookingRequest struct {
	GuestID     int    `json:"guest_id" validate:"required,min=1"`
	RoomID      int    `json:"room_id" validate:"required,min=1"`
	BookingDate string `json:"booking_date" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
	CreatedBy   string `json:"created_by" validate:"required,min=1,max=100"`
}

type PatchBookingRequest struct {
	BookingDate   *string `json:"booking_date,omitempty"`
	StartTime     *string `json:"start_time,omitempty"`
	EndTime       *string `json:"end_time,omitempty"`
	BookingStatus *string `json:"booking_status,omitempty" validate:"omitempty,oneof=active cancelled"`
	IsActive      *bool   `json:"is_active,omitempty"`
	UpdatedBy     string  `json:"updated_by" validate:"required,min=1,max=100"`
}

type BookingResponse struct {
	BookingUUID   string `json:"booking_uuid"`
	GuestID       int    `json:"guest_id"`
	RoomID        int    `json:"room_id"`
	BookingDate   string `json:"booking_date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	BookingStatus string `json:"booking_status"`
	IsActive      bool   `json:"is_active"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type BookingListResponse struct {
	Bookings []BookingResponse `json:"bookings"`
	Total    int64             `json:"total"`
}

type APIResponse struct {
	Status    string      `json:"status"`
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    []Error     `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp"`
}

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
