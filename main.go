package main

import (
	"fmt"
	"job-application-automation/database"
	"job-application-automation/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Connecting to the database...")
	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	fmt.Println("Running migrations...")
	if err := database.Migrate(); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	e := echo.New()
	routes.InitRoutes(e)

	fmt.Println("Starting server on port 8082...")
	e.Logger.Fatal(e.Start(":8082"))
}
