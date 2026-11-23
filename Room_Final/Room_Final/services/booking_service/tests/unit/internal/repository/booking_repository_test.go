package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"booking_system/booking_service/internal/loggers"
	"booking_system/booking_service/internal/models"
	"booking_system/booking_service/internal/repository"
	"booking_system/booking_service/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BookingRepositoryTestSuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository repository.BookingRepositoryInterface
}

func (suite *BookingRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	suite.NoError(err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	suite.NoError(err)

	dbMock := &database.Db{
		Gorm:  gormDB,
		SqlDb: db,
	}

	logger := loggers.NewLogger(loggers.LogConfig{
		Level:       "info",
		LogDir:      "./logs",
		FileName:    "test.log",
		ServiceName: "test-service",
	})

	suite.DB = gormDB
	suite.mock = mock
	suite.repository = repository.NewBookingRepository(dbMock, logger)
}

// ---- CreateBooking ----

func (suite *BookingRepositoryTestSuite) TestCreateBooking_Success() {
	booking := &models.Booking{
		BookingUUID:      "test-uuid-001",
		UserID:           1,
		RoomID:           1,
		BookingDate:      time.Now().Add(24 * time.Hour),
		StartTime:        "09:00",
		EndTime:          "10:00",
		NumberOfEmployee: 5,
		BookingStatus:    "active",
		IsActive:         true,
		CreatedBy:        "admin",
		UpdatedBy:        "admin",
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bookings"`)).
		WillReturnRows(sqlmock.NewRows([]string{"booking_id"}).AddRow(1))
	suite.mock.ExpectCommit()

	err := suite.repository.CreateBooking(context.Background(), booking)

	assert.Nil(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestCreateBooking_DBError() {
	booking := &models.Booking{
		BookingUUID: "test-uuid-002",
		UserID:      1,
		RoomID:      1,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bookings"`)).
		WillReturnError(sqlmock.ErrCancelled)
	suite.mock.ExpectRollback()

	err := suite.repository.CreateBooking(context.Background(), booking)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Message, "Failed to create booking")
}

// ---- GetAllBookings ----

func (suite *BookingRepositoryTestSuite) TestGetAllBookings_Success() {
	now := time.Now()
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(countRows)

	dataRows := sqlmock.NewRows([]string{
		"booking_id", "booking_uuid", "user_id", "room_id",
		"booking_date", "start_time", "end_time", "number_of_employee",
		"booking_status", "is_active", "created_by", "updated_by",
		"created_at", "updated_at",
	}).
		AddRow(1, "uuid-1", 1, 1, now, "09:00", "10:00", 5, "active", true, "admin", "admin", now, now).
		AddRow(2, "uuid-2", 2, 2, now, "11:00", "12:00", 3, "active", true, "admin", "admin", now, now)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(dataRows)

	bookings, count, err := suite.repository.GetAllBookings(context.Background())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), int64(2), count)
	assert.Len(suite.T(), bookings, 2)
}

func (suite *BookingRepositoryTestSuite) TestGetAllBookings_CountError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnError(sqlmock.ErrCancelled)

	bookings, count, err := suite.repository.GetAllBookings(context.Background())

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
	assert.Nil(suite.T(), bookings)
}

func (suite *BookingRepositoryTestSuite) TestGetAllBookings_FetchError() {
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(countRows)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(sqlmock.ErrCancelled)

	bookings, count, err := suite.repository.GetAllBookings(context.Background())

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), int64(0), count)
	assert.Nil(suite.T(), bookings)
}

// ---- GetBookingByUUID ----

func (suite *BookingRepositoryTestSuite) TestGetBookingByUUID_Success() {
	now := time.Now()
	uuid := "test-uuid-001"

	rows := sqlmock.NewRows([]string{
		"booking_id", "booking_uuid", "user_id", "room_id",
		"booking_date", "start_time", "end_time", "number_of_employee",
		"booking_status", "is_active", "created_by", "updated_by",
		"created_at", "updated_at",
	}).AddRow(1, uuid, 1, 1, now, "09:00", "10:00", 5, "active", true, "admin", "admin", now, now)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(rows)

	booking, err := suite.repository.GetBookingByUUID(context.Background(), uuid)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), booking)
	assert.Equal(suite.T(), uuid, booking.BookingUUID)
}

func (suite *BookingRepositoryTestSuite) TestGetBookingByUUID_NotFound() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(gorm.ErrRecordNotFound)

	booking, err := suite.repository.GetBookingByUUID(context.Background(), "non-existent")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), booking)
	assert.Contains(suite.T(), err.Message, "not found")
}

func (suite *BookingRepositoryTestSuite) TestGetBookingByUUID_DBError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(sqlmock.ErrCancelled)

	booking, err := suite.repository.GetBookingByUUID(context.Background(), "test-uuid")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), booking)
}

// ---- PatchBooking ----

func (suite *BookingRepositoryTestSuite) TestPatchBooking_Success() {
	updateData := map[string]interface{}{
		"start_time": "10:00",
		"updated_by": "admin",
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.PatchBooking(context.Background(), "test-uuid-001", updateData)

	assert.Nil(suite.T(), err)
}

func (suite *BookingRepositoryTestSuite) TestPatchBooking_NotFound() {
	updateData := map[string]interface{}{"updated_by": "admin"}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings"`)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	suite.mock.ExpectCommit()

	err := suite.repository.PatchBooking(context.Background(), "non-existent", updateData)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Message, "not found")
}

func (suite *BookingRepositoryTestSuite) TestPatchBooking_DBError() {
	updateData := map[string]interface{}{"updated_by": "admin"}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "bookings"`)).
		WillReturnError(sqlmock.ErrCancelled)
	suite.mock.ExpectRollback()

	err := suite.repository.PatchBooking(context.Background(), "test-uuid-001", updateData)

	assert.NotNil(suite.T(), err)
}

// ---- UserExists ----

func (suite *BookingRepositoryTestSuite) TestUserExists_True() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.UserExists(context.Background(), 1)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *BookingRepositoryTestSuite) TestUserExists_False() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.UserExists(context.Background(), 999)

	assert.Nil(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *BookingRepositoryTestSuite) TestUserExists_DBError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnError(sqlmock.ErrCancelled)

	exists, err := suite.repository.UserExists(context.Background(), 1)

	assert.NotNil(suite.T(), err)
	assert.False(suite.T(), exists)
}

// ---- RoomExists ----

func (suite *BookingRepositoryTestSuite) TestRoomExists_True() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.RoomExists(context.Background(), 1)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *BookingRepositoryTestSuite) TestRoomExists_False() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.RoomExists(context.Background(), 999)

	assert.Nil(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *BookingRepositoryTestSuite) TestRoomExists_DBError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnError(sqlmock.ErrCancelled)

	exists, err := suite.repository.RoomExists(context.Background(), 1)

	assert.NotNil(suite.T(), err)
	assert.False(suite.T(), exists)
}

// ---- CheckRoomAvailability ----

func (suite *BookingRepositoryTestSuite) TestCheckRoomAvailability_Available() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	available, err := suite.repository.CheckRoomAvailability(
		context.Background(), 1,
		time.Now().Add(24*time.Hour),
		"09:00", "10:00", "",
	)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), available)
}

func (suite *BookingRepositoryTestSuite) TestCheckRoomAvailability_NotAvailable() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	available, err := suite.repository.CheckRoomAvailability(
		context.Background(), 1,
		time.Now().Add(24*time.Hour),
		"09:00", "10:00", "",
	)

	assert.Nil(suite.T(), err)
	assert.False(suite.T(), available)
}

func (suite *BookingRepositoryTestSuite) TestCheckRoomAvailability_DBError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnError(sqlmock.ErrCancelled)

	available, err := suite.repository.CheckRoomAvailability(
		context.Background(), 1,
		time.Now().Add(24*time.Hour),
		"09:00", "10:00", "",
	)

	assert.NotNil(suite.T(), err)
	assert.False(suite.T(), available)
}

func TestBookingRepositorySuite(t *testing.T) {
	suite.Run(t, new(BookingRepositoryTestSuite))
}

