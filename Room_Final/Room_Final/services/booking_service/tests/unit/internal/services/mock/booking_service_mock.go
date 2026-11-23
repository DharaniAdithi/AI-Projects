package mock

import (
	"context"

	"booking_system/booking_service/internal/dtos"

	"github.com/stretchr/testify/mock"
)

// BookingServiceMock mocks BookingServiceInterface
type BookingServiceMock struct {
	mock.Mock
}

func (m *BookingServiceMock) CreateBooking(ctx context.Context, req dtos.CreateBookingRequest) *dtos.APIResponse {
	args := m.Called(ctx, req)
	return args.Get(0).(*dtos.APIResponse)
}

func (m *BookingServiceMock) GetAllBookings(ctx context.Context) *dtos.APIResponse {
	args := m.Called(ctx)
	return args.Get(0).(*dtos.APIResponse)
}

func (m *BookingServiceMock) GetBookingByUUID(ctx context.Context, bookingUUID string) *dtos.APIResponse {
	args := m.Called(ctx, bookingUUID)
	return args.Get(0).(*dtos.APIResponse)
}

func (m *BookingServiceMock) PatchBooking(ctx context.Context, bookingUUID string, req dtos.PatchBookingRequest) *dtos.APIResponse {
	args := m.Called(ctx, bookingUUID, req)
	return args.Get(0).(*dtos.APIResponse)
}

