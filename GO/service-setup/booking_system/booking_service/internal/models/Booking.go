package models

import (
	"time"
)

// User represents a customer
type User struct {
	UserID    int       `gorm:"primaryKey;autoIncrement" json:"user_id"`
	UserName  string    `gorm:"column:user_name;type:varchar(100);not null" json:"user_name"`
	Email     string    `gorm:"column:email;type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"column:phone;type:varchar(20)" json:"phone"`
	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Room represents a room with capacity and pricing
type Room struct {
	RoomID      int       `gorm:"primaryKey;autoIncrement" json:"room_id"`
	RoomName    string    `gorm:"column:room_name;type:varchar(100);not null" json:"room_name"`
	Capacity    int       `gorm:"column:capacity;not null" json:"capacity"`
	PricePerDay int       `gorm:"column:price_per_day;not null" json:"price_per_day"` // in INR
	IsActive    bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Room) TableName() string {
	return "rooms"
}

// Booking represents a room booking
type Booking struct {
	BookingID     uint      `gorm:"primaryKey;autoIncrement" json:"booking_id"`
	BookingUUID   string    `gorm:"column:booking_uuid;type:varchar(36);uniqueIndex;not null" json:"booking_uuid"`
	UserID        int       `gorm:"column:user_id;not null;index" json:"user_id"`
	RoomID        int       `gorm:"column:room_id;not null;index" json:"room_id"`
	BookingDate   time.Time `gorm:"column:booking_date;not null;index" json:"booking_date"`
	StartTime     string    `gorm:"column:start_time;type:varchar(10);not null" json:"start_time"`
	EndTime       string    `gorm:"column:end_time;type:varchar(10);not null" json:"end_time"`
	PricePerDay   int       `gorm:"column:price_per_day;not null" json:"price_per_day"` // in INR, captured at booking time
	BookingStatus string    `gorm:"column:booking_status;type:varchar(100);default:'active'" json:"booking_status"`
	IsActive      bool      `gorm:"column:is_active;default:true;index" json:"is_active"` // soft delete
	CreatedBy     string    `gorm:"column:created_by;type:varchar(100);not null" json:"created_by"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedBy     string    `gorm:"column:updated_by;type:varchar(100);not null" json:"updated_by"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Booking) TableName() string {
	return "bookings"
}
