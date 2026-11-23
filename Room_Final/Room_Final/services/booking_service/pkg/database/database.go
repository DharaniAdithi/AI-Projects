// package database

// import (
// 	"database/sql"
// 	"strings"

// 	_ "github.com/mattn/go-sqlite3"
// 	_ "github.com/microsoft/go-mssqldb"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/driver/sqlserver"
// 	"gorm.io/gorm"
// )

// type DBConnector interface {
// 	EstablishConnection(dbURL string) (*Db, error)
// }

// type Db struct {
// 	Gorm  *gorm.DB
// 	SqlDb *sql.DB
// }

// type DBService struct{}

// func (s *DBService) EstablishConnection(dbURL string) (*Db, error) {
// 	var dialector gorm.Dialector

// 	if strings.HasPrefix(dbURL, "sqlite://") {
// 		dbPath := strings.TrimPrefix(dbURL, "sqlite://")
// 		dialector = sqlite.Open(dbPath)
// 	} else if strings.HasPrefix(dbURL, "sqlserver://") {
// 		dialector = sqlserver.Open(dbURL)
// 	} else {
// 		dialector = sqlserver.Open(dbURL)
// 	}

// 	db, err := gorm.Open(dialector, &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return &Db{
// 		Gorm:  db,
// 		SqlDb: sqlDB,
// 	}, nil
// }

package database

import (
	"database/sql"
	"fmt"
	"strings"

	
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DBConnector interface {
	EstablishConnection(dbURL string) (*Db, error)
}

type Db struct {
	Gorm  *gorm.DB
	SqlDb *sql.DB
}

type DBService struct{}  //group related db methods

func (s *DBService) EstablishConnection(dbURL string) (*Db, error) {
	var dialector gorm.Dialector   //connect to sql ,postgress

	switch {
	case strings.HasPrefix(dbURL, "sqlite://"):
		dbPath := strings.TrimPrefix(dbURL, "sqlite://")
		dialector = sqlite.Open(dbPath)

	case strings.HasPrefix(dbURL, "sqlserver://"):
		dialector = sqlserver.Open(dbURL)

	case strings.HasPrefix(dbURL, "postgres://"),
		strings.HasPrefix(dbURL, "postgresql://"):
		dialector = postgres.Open(dbURL)

	default:
		return nil, fmt.Errorf("unsupported database URL: %s", dbURL)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})  //tries to open database connection
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {  //alive or not
		return nil, err
	}

	return &Db{
		Gorm:  db,
		SqlDb: sqlDB,
	}, nil
}
