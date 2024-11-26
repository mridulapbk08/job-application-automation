package main

import (
	"fmt"
	"job-application-automation/database"
	"job-application-automation/models"
	"job-application-automation/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Connecting to the database...")
	database.Connect()

	fmt.Println("Running auto-migrations...")
	err := database.DB.AutoMigrate(&models.Job{}, &models.Tracker{})
	if err != nil {
		fmt.Printf("Failed to auto-migrate tables: %v\n", err)
		return
	}
	fmt.Println("Auto-migrations completed.")

	e := echo.New()
	routes.InitRoutes(e)

	fmt.Println("Starting the server on port 8082...")
	e.Logger.Fatal(e.Start(":8082"))
}
