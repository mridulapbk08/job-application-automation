package routes

import (
	"github.com/labstack/echo/v4"
	"job-application-automation/controllers"
)

func InitRoutes(e *echo.Echo) {
	e.POST("/apply", controllers.ApplyForJob)
}
