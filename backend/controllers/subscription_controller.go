package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

type SubscriptionRequest struct {
	PlanID          uint   `json:"plan_id" validate:"required"`
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
}

// @Summary Create subscription
// @Description Subscribe to a creator's plan
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body SubscriptionRequest true "Subscription data"
// @Success 201 {object} models.Subscription
// @Router /subscriptions [post]
func CreateSubscription(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}

	var req SubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	db := c.Get("db").(*gorm.DB)
	
	// Verificar que el plan existe
	var plan models.Plan
	if err := db.First(&plan, req.PlanID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Plan no encontrado"})
	}

	// TODO: Integración con Stripe aquí
	
	subscription := models.Subscription{
		UserID:                 user.ID,
		CreatorID:              plan.CreatorID,
		PlanID:                 plan.ID,
		StripeSubscriptionID:   "temp_stripe_id", // Reemplazar con ID real de Stripe
		Status:                 "active",
	}

	if err := db.Create(&subscription).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear suscripción"})
	}

	return c.JSON(http.StatusCreated, subscription)
}

// @Summary Get user subscriptions
// @Description List all active subscriptions for current user
// @Tags Subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Router /subscriptions [get]
func GetUserSubscriptions(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}

	db := c.Get("db").(*gorm.DB)
	
	var subscriptions []models.Subscription
	if err := db.Where("user_id = ?", user.ID).Find(&subscriptions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al obtener suscripciones"})
	}

	return c.JSON(http.StatusOK, subscriptions)
}
// GetSubscriptionDetails obtiene detalles de una suscripción
// @Summary Get subscription details
// @Description Get details of a specific subscription
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Router /subscriptions/{id} [get]
func GetSubscriptionDetails(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}

	subscriptionID := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	
	var subscription models.Subscription
	if err := db.Preload("Plan").Preload("Creator").
		Where("id = ? AND user_id = ?", subscriptionID, user.ID).
		First(&subscription).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Suscripción no encontrada"})
	}

	return c.JSON(http.StatusOK, subscription)
}

// CancelSubscription cancela una suscripción
// @Summary Cancel subscription
// @Description Cancel an active subscription
// @Tags Subscriptions
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]interface{}
// @Router /subscriptions/{id} [delete]
func CancelSubscription(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}

	subscriptionID := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	
	// Verificar que la suscripción existe y pertenece al usuario
	var subscription models.Subscription
	if err := db.Where("id = ? AND user_id = ?", subscriptionID, user.ID).
		First(&subscription).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Suscripción no encontrada"})
	}

	// TODO: Integración con Stripe para cancelar suscripción
	
	// Actualizar estado en la base de datos
	if err := db.Model(&subscription).
		Updates(map[string]interface{}{
			"status":    "canceled",
			"end_date":  time.Now(),
		}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al cancelar suscripción"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Suscripción cancelada exitosamente",
		"end_date": subscription.EndDate,
	})
}