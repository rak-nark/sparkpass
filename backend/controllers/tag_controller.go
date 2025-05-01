package controllers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/rak-nark/sparkpass/models"
    "gorm.io/gorm"
    "github.com/gosimple/slug"
)

// GetAllTags - Listar tags disponibles
// @Summary List all tags
// @Description Get all available content tags
// @Tags Tags
// @Produce json
// @Success 200 {array} models.Tag
// @Router /tags [get]
func GetAllTags(c echo.Context) error {
    db := c.Get("db").(*gorm.DB)
    
    var tags []models.Tag
    if err := db.Find(&tags).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Error al obtener etiquetas",
        })
    }
    
    return c.JSON(http.StatusOK, tags)
}

type CreateTagRequest struct {
    Name string `json:"name" validate:"required,max=50"`
}

// CreateTag - Crear nueva etiqueta (accesible para creadores)
// @Summary Create a tag
// @Description Create a new content tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param request body CreateTagRequest true "Tag data"
// @Success 201 {object} models.Tag
// @Router /tags [post]
func CreateTag(c echo.Context) error {
    // Verificar que el usuario es creador
    user, ok := c.Get("user").(*models.User)
    if !ok || !user.IsCreator {
        return c.JSON(http.StatusForbidden, map[string]string{
            "error": "Solo creadores pueden crear etiquetas",
        })
    }

    var req CreateTagRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Datos inválidos",
        })
    }

    db := c.Get("db").(*gorm.DB)
    
    // Verificar si el tag ya existe
    var existingTag models.Tag
    if err := db.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
        return c.JSON(http.StatusConflict, map[string]string{
            "error": "La etiqueta ya existe",
        })
    }

    // Crear nuevo tag
    newTag := models.Tag{
        Name: req.Name,
        Slug: slug.Make(req.Name),
    }

    if err := db.Create(&newTag).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Error al crear etiqueta",
        })
    }

    return c.JSON(http.StatusCreated, newTag)
}