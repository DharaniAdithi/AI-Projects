package dtos

type CreateBookingRequest struct {
	UserID           int    `json:"user_id" validate:"required,min=1"`
	RoomID           int    `json:"room_id" validate:"required,min=1"`
	BookingDate      string `json:"booking_date" validate:"required"`
	StartTime        string `json:"start_time" validate:"required"`
	EndTime          string `json:"end_time" validate:"required"`
	NumberOfEmployee int    `json:"number_of_employee" validate:"required,min=1"`
	CreatedBy        string `json:"created_by" validate:"required,min=1,max=100"`
}

type PatchBookingRequest struct {
	BookingDate      *string `json:"booking_date,omitempty"`
	StartTime        *string `json:"start_time,omitempty"`
	EndTime          *string `json:"end_time,omitempty"`
	RoomID           *int    `json:"room_id,omitempty" validate:"omitempty,min=1"`
	NumberOfEmployee *int    `json:"number_of_employee,omitempty" validate:"omitempty,min=1"`
	IsActive         *bool   `json:"is_active,omitempty"`
	UpdatedBy        string  `json:"updated_by" validate:"required,min=1,max=100"`
}

type BookingResponse struct {
	BookingUUID   string `json:"booking_uuid"`
	UserID        int    `json:"user_id"`
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

// type UserResponse struct {
// 	UserID   int    `json:"user_id"`
// 	UserName string `json:"user_name"`
// 	Email    string `json:"email"`
// 	Phone    string `json:"phone"`
// 	IsActive bool   `json:"is_active"`
// }

// type RoomResponse struct {
// 	RoomID   int    `json:"room_id"`
// 	RoomName string `json:"room_name"`
// 	Capacity int    `json:"capacity"`
// 	IsActive bool   `json:"is_active"`
// }

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
