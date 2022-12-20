package model

import (
	"database/sql"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gDB     *gorm.DB
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	err     error
)

func setup() {
	mockDB, mockSQL, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	dialector := postgres.New(postgres.Config{
		DriverName: "postgres",
		DSN:        "sql_mock_db_0",
		Conn:       mockDB,
	})
	gDB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func teardown() {
	mockDB.Close()
}

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}
