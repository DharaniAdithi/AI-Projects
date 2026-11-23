package database_test

import (
	"testing"

	"booking_system/booking_service/pkg/database"

	"github.com/stretchr/testify/assert"
)

func TestEstablishConnection_UnsupportedDriver(t *testing.T) {
	dbService := &database.DBService{}

	db, err := dbService.EstablishConnection("mysql://user:pass@localhost/testdb")

	assert.Nil(t, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database URL")
}

func TestEstablishConnection_EmptyURL(t *testing.T) {
	dbService := &database.DBService{}

	db, err := dbService.EstablishConnection("")

	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestEstablishConnection_InvalidPostgresURL(t *testing.T) {
	dbService := &database.DBService{}

	// Postgres URL with unreachable host — fails at Ping()
	db, err := dbService.EstablishConnection("postgres://invalid:invalid@localhost:9999/nodb")

	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestEstablishConnection_InvalidSQLServerURL(t *testing.T) {
	dbService := &database.DBService{}

	db, err := dbService.EstablishConnection("sqlserver://invalid:invalid@localhost:9999?database=nodb")

	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestEstablishConnection_PostgresPrefix_Recognized(t *testing.T) {
	dbService := &database.DBService{}

	// "postgresql://" prefix should be recognized (not return "unsupported" error)
	db, err := dbService.EstablishConnection("postgresql://invalid:invalid@localhost:9999/nodb")

	// It will fail to connect, but NOT with "unsupported database URL"
	assert.Nil(t, db)
	assert.Error(t, err)
	assert.NotContains(t, err.Error(), "unsupported database URL")
}
