package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/controllers"
	"github.com/rak-nark/sparkpass/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	authMiddleware := middleware.AuthMiddleware(db)

	// Auth (Públicas)
	e.POST("/api/register", controllers.Register)
	e.POST("/api/login", controllers.Login)

	// Grupo general protegido
	api := e.Group("/api")
	api.Use(authMiddleware)
	
	// Rutas de Usuario
	user := api.Group("/users")
	{
		user.GET("/me", controllers.GetUserProfile)
		user.PUT("/me", controllers.UpdateUserProfile)
	}

	// Rutas de Creadores
	creators := api.Group("/creators")
	{
		creators.POST("", controllers.BecomeCreator)
		creators.GET("/:id", controllers.GetCreatorProfile)
		creators.GET("/me/plans", controllers.GetCreatorPlans)
		creators.POST("/plans", controllers.CreatePlan)
	}

	// Rutas de Contenido
	content := api.Group("/content")
	{
		content.GET("", controllers.GetContent)
		content.GET("/:id", controllers.GetContentByID)
		content.POST("", controllers.CreateContent)
		content.POST("/:id/tags", controllers.AddTagsToContent)
		content.PUT("/:id", controllers.UpdateContent)
		content.DELETE("/:id", controllers.DeleteContent)
		content.GET("/search", controllers.SearchContent)
	}

	// Rutas de Suscripciones
	subscriptions := api.Group("/subscriptions")
	{
		subscriptions.GET("", controllers.GetUserSubscriptions)
		subscriptions.POST("", controllers.CreateSubscription)
		subscriptions.GET("/:id", controllers.GetSubscriptionDetails)
		subscriptions.DELETE("/:id", controllers.CancelSubscription)
	}

	// Rutas de Pagos
	payments := api.Group("/payments")
	{
		payments.GET("", controllers.GetPaymentHistory)
		payments.GET("/:id", controllers.GetPaymentDetails)
	}

	// Rutas de Tags
	tags := api.Group("/tags")
	{
		tags.GET("", controllers.GetAllTags)
		tags.POST("", controllers.CreateTag)
	}
}