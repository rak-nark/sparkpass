package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/controllers"
	"github.com/rak-nark/sparkpass/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Inicializar middleware con la instancia de DB
	authMiddleware := middleware.AuthMiddleware(db)

	// Auth (Públicas)
	e.POST("/api/register", controllers.Register)
	e.POST("/api/login", controllers.Login)

	// Grupo general protegido
	protected := e.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/content", controllers.GetContent)
		protected.POST("/creators", controllers.BecomeCreator)
	}
}