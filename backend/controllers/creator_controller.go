package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

type CreatorRequest struct {
	Bio string `json:"bio" validate:"required,max=500"`
}
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
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
// GetCreatorProfile obtiene el perfil público de un creador
// @Summary Get creator profile
// @Description Get public profile of a creator
// @Tags Creators
// @Produce json
// @Param id path int true "Creator ID"
// @Success 200 {object} map[string]interface{}
// @Router /creators/{id} [get]
func GetCreatorProfile(c echo.Context) error {
	creatorID := c.Param("id")
	
	db := c.Get("db").(*gorm.DB)
	
	var creator struct {
		models.Creator
		User struct {
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"user"`
		ContentCount int64 `json:"content_count"`
	}
	
	// Obtener información básica del creador
	if err := db.Table("creators").
		Select("creators.*, users.email, users.avatar_url").
		Joins("LEFT JOIN users ON users.id = creators.user_id").
		Where("creators.user_id = ?", creatorID).
		Scan(&creator).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Creador no encontrado"})
	}
	
	// Contar contenido del creador
	db.Model(&models.PremiumContent{}).
		Where("creator_id = ?", creatorID).
		Count(&creator.ContentCount)
	
	return c.JSON(http.StatusOK, creator)
}

type PlanRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description" validate:"required,max=500"`
	BasePrice   float64 `json:"base_price" validate:"required,min=0"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// CreatePlan crea un nuevo plan de suscripción
// @Summary Create subscription plan
// @Description Create a new subscription plan (for creators)
// @Tags Creators
// @Accept json
// @Produce json
// @Param request body PlanRequest true "Plan data"
// @Success 201 {object} models.Plan
// @Router /creators/plans [post]
func CreatePlan(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || !user.IsCreator {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Solo creadores pueden crear planes"})
	}

	var req PlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	db := c.Get("db").(*gorm.DB)
	
	// TODO: Integración con Stripe para crear precio
	stripePriceID := "price_temp_" + generateRandomString(10)
	
	plan := models.Plan{
		CreatorID:      user.ID,
		Name:          req.Name,
		Description:   req.Description,
		BasePrice:     req.BasePrice,
		Price:         req.Price,
		StripePriceID: stripePriceID,
		IsActive:      true,
	}

	if err := db.Create(&plan).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear plan"})
	}

	return c.JSON(http.StatusCreated, plan)
}

// GetCreatorPlans lista los planes de un creador
// @Summary Get creator plans
// @Description List all plans from current creator
// @Tags Creators
// @Produce json
// @Success 200 {array} models.Plan
// @Router /creators/me/plans [get]
func GetCreatorPlans(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || !user.IsCreator {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Solo creadores pueden ver sus planes"})
	}

	db := c.Get("db").(*gorm.DB)
	
	var plans []models.Plan
	if err := db.Where("creator_id = ?", user.ID).Find(&plans).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al obtener planes"})
	}

	return c.JSON(http.StatusOK, plans)
}