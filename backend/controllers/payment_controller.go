package controllers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

// @Summary Get payment history
// @Description List all payments for current user
// @Tags Payments
// @Produce json
// @Success 200 {array} models.Payment
// @Router /payments [get]
func GetPaymentHistory(c echo.Context) error {
    user, ok := c.Get("user").(*models.User)
    if !ok {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "error": "Usuario no autenticado",
        })
    }

    db := c.Get("db").(*gorm.DB)
    
    // Debug: Imprime el ID del usuario
    log.Printf("Consultando pagos para el usuario ID: %d", user.ID)

    var payments []models.Payment
    err := db.
        Debug(). // Muestra el SQL en consola
        Preload("Subscription").
        Joins("JOIN subscriptions ON subscriptions.id = payments.subscription_id").
        Where("subscriptions.user_id = ?", user.ID).
        Order("payments.created_at DESC").
        Find(&payments).Error

    // Debug: Imprime los resultados
    log.Printf("Encontrados %d pagos", len(payments))
    for _, p := range payments {
        log.Printf("Pago ID: %d, Amount: %.2f, SubscriptionID: %d", p.ID, p.Amount, p.SubscriptionID)
    }

    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Error al obtener historial de pagos: " + err.Error(),
        })
    }

    return c.JSON(http.StatusOK, payments)
}

// @Summary Get payment details
// @Description Get details of a specific payment
// @Tags Payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} models.Payment
// @Router /payments/{id} [get]
func GetPaymentDetails(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	paymentID := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	
	var payment models.Payment
	if err := db.Preload("Subscription").
		Joins("JOIN subscriptions ON subscriptions.id = payments.subscription_id").
		Where("payments.id = ? AND subscriptions.user_id = ?", paymentID, user.ID).
		First(&payment).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Pago no encontrado",
		})
	}

	return c.JSON(http.StatusOK, payment)
}