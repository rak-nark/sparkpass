package controllers

import (
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

type CreatorRequest struct {
	Bio string `json:"bio" validate:"required,max=500"`
}
func BecomeCreator(c echo.Context) error {
	// Obtener usuario del contexto
	user, ok := c.Get("user").(*models.User)
	if !ok || user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "No autorizado",
			"message": "Debe iniciar sesión para realizar esta acción",
		})
	}

	// Validar que el usuario no sea ya creador
	if user.IsCreator {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error":   "Usuario ya es creador",
			"message": "El usuario ya tiene perfil de creador",
		})
	}

	// Parsear request
	var req CreatorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Datos inválidos",
			"message": err.Error(),
		})
	}

	// Validar datos
	if err := c.Validate(req); err != nil {
		// Mejorar mensajes de error de validación
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages[i] = "El campo " + e.Field() + " es requerido"
			case "max":
				errorMessages[i] = "El campo " + e.Field() + " debe tener máximo " + e.Param() + " caracteres"
			default:
				errorMessages[i] = "Error en el campo " + e.Field()
			}
		}
		
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validación fallida",
			"details": errorMessages,
		})
	}

	// Obtener DB del contexto
	db := c.Get("db").(*gorm.DB)

	// Transacción para asegurar consistencia
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Actualizar usuario como creador
		if err := tx.Model(user).Update("is_creator", true).Error; err != nil {
			return err
		}

		// 2. Crear perfil de creador
		creator := models.Creator{
			UserID:          user.ID,
			Bio:             req.Bio,
			StripeAccountID: "pending_activation", // Se actualizará con Stripe Connect
		}

		if err := tx.Create(&creator).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Error al crear perfil",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Perfil de creador creado exitosamente",
		"next_steps": []string{
			"Complete su configuración de Stripe para recibir pagos",
			"Configure sus planes de suscripción",
		},
	})
}