package database

import (
	"fmt"
	"job-application-automation/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "root:Root@1234#@tcp(localhost:3306)/job_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	return nil
}

func Migrate() error {
	if err := DB.AutoMigrate(&models.Job{}, &models.Tracker{}); err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %v", err)
	}
	return nil
}
