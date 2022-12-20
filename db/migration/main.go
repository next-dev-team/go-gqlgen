package main

import (
	"fmt"
	"go-graph/db"
	"go-graph/db/model"
	"go-graph/pkg/config"
	"log"
)

func main() {
	if err := config.InitDefaultGormConfig(); err != nil {
		panic(err)
	}
	conn := db.GetConnection()
	if err := conn.AutoMigrate(
		&model.Todo{},
	); err != nil {
		panic(fmt.Errorf("automatically migrate database failed %v", err))
	}
	log.Println("migrate tables created...")
}
