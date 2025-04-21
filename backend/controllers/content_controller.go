package controllers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/config"
	"github.com/rak-nark/sparkpass/models"
)

// @Summary Get all content
// @Description List all premium content
// @Tags Content
// @Produce json
// @Success 200 {array} models.PremiumContent
// @Router /content [get]
func GetContent(c echo.Context) error {
	var content []models.PremiumContent
	config.DB.Find(&content)
	return c.JSON(http.StatusOK, content)
}