package db

import (
    "test_go_project/internal/entity"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
	"fmt"
)

func Init(host string, user string, password string, name string) *gorm.DB {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, user, password, name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Failed to connect to db: ", err)
    }

    err = db.AutoMigrate(&entity.RefreshRecord{})

    if err != nil {
        log.Fatal("Failed to connect to db: ", err)
    }
	
    return db
}