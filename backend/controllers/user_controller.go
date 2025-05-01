package controllers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

// @Summary Get user profile
// @Description Get current authenticated user's profile
// @Tags Users
// @Produce json
// @Success 200 {object} models.User
// @Router /users/me [get]
func GetUserProfile(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}
	
	// Ocultar campos sensibles
	user.PasswordHash = ""
	return c.JSON(http.StatusOK, user)
}

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

// @Summary Update user profile
// @Description Update current user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Update data"
// @Success 200 {object} models.User
// @Router /users/me [put]
func UpdateUserProfile(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	db := c.Get("db").(*gorm.DB)
	
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	if err := db.Model(user).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al actualizar perfil"})
	}

	return c.JSON(http.StatusOK, user)
}