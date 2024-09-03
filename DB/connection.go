package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func NewDBConnection() *gorm.DB {

	url := "root:root@tcp(localhost:3307)/fooddelight?charset=utf8mb4&parseTime=true"

	db, err := gorm.Open("mysql", url)
	if err != nil {
		fmt.Println("Error in Database Connectivity...")
		return nil
	}

	// Sets connect restrictions
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	db.LogMode(true)

	return db
}
