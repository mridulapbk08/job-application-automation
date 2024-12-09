package routes

import (
	"job-application-automation/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.POST("/apply", controllers.ApplyForJob)
}
