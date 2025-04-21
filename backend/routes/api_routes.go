package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/controllers"
	"github.com/rak-nark/sparkpass/middleware"
)

func SetupRoutes(e *echo.Echo) {
	// Auth
	e.POST("/api/register", controllers.Register)
	e.POST("/api/login", controllers.Login)

	// Rutas protegidas
	contentGroup := e.Group("/api/content")
	contentGroup.Use(middleware.AuthMiddleware)
	{
		contentGroup.GET("", controllers.GetContent)
		// Agregar más rutas protegidas aquí
	}
}