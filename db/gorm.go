package db

import (
	"fmt"
	"go-graph/pkg/config"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	conf := config.GetGormConfig()
	fmt.Println("dsn:", conf.GetDsn())
	db, err := gorm.Open(postgres.Open(conf.GetDsn()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error opening database: %v", err))
	}
	//! Connection Pool
	sqlDB, _ := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(conf.GetMaxIdleConn())
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(conf.GetMaxOpenConn())
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(conf.GetMaxLifetimeConn())
	log.Info().Msg("gorm initialized")
	return db
}

func GetTx() (*gorm.DB, error) {
	conn := GetConnection()
	tx := conn.Begin()
	if err := tx.Error; err != nil {
		log.Err(err).Msg("start transaction error")
		return nil, err
	}
	return tx, nil
}
