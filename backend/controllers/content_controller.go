package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"gorm.io/gorm"
)

type ContentRequest struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
	ContentType string  `json:"content_type" validate:"required,oneof=video image podcast document"`
	IsLocked    bool    `json:"is_locked"`
	S3Key       string  `json:"s3Key"`
    Slug        string  `json:"slug"` 
}

// @Summary Create content
// @Description Create new premium content (for creators)
// @Tags Content
// @Accept json
// @Produce json
// @Param request body ContentRequest true "Content data"
// @Success 201 {object} models.PremiumContent
// @Router /content [post]
func CreateContent(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || !user.IsCreator {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Solo creadores pueden publicar contenido"})
	}

	var req ContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	db := c.Get("db").(*gorm.DB)
	
	content := models.PremiumContent{
		CreatorID:   user.ID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		ContentType: req.ContentType,
		IsLocked:    req.IsLocked,
		Slug:        slug.Make(req.Title),
	}

	if err := db.Create(&content).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear contenido"})
	}

	return c.JSON(http.StatusCreated, content)
}
// @Summary Get all content
// @Description List all premium content with optional filters
// @Tags Content
// @Produce json
// @Param is_locked query bool false "Filter by locked status"
// @Param content_type query string false "Filter by content type (video,image,podcast,document)"
// @Param creator_id query int false "Filter by creator ID"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param limit query int false "Results per page (default 20, max 100)"
// @Param page query int false "Page number (default 1)"
// @Success 200 {array} models.PremiumContent
// @Router /content [get]
func GetContent(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	// Leer parámetros manualmente
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	isLockedParam := c.QueryParam("is_locked")
	var isLocked *bool
	if isLockedParam != "" {
		val := isLockedParam == "true"
		isLocked = &val
	}

	contentType := c.QueryParam("content_type")
	creatorID, _ := strconv.Atoi(c.QueryParam("creator_id"))
	minPrice, _ := strconv.ParseFloat(c.QueryParam("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.QueryParam("max_price"), 64)

	// Construir consulta
	query := db.Model(&models.PremiumContent{}).
		Preload("Creator").
		Limit(limit).
		Offset((page - 1) * limit).
		Order("created_at DESC")

	if isLocked != nil {
		query = query.Where("is_locked = ?", *isLocked)
	}
	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}
	if creatorID != 0 {
		query = query.Where("creator_id = ?", creatorID)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	var content []models.PremiumContent
	if err := query.Find(&content).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al obtener contenido",
		})
	}

	return c.JSON(http.StatusOK, content)
}


// @Summary Get content by ID
// @Description Get specific content details
// @Tags Content
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} models.PremiumContent
// @Router /content/{id} [get]
func GetContentByID(c echo.Context) error {
	id := c.Param("id")
	
	var content models.PremiumContent
	db := c.Get("db").(*gorm.DB)
	
	if err := db.First(&content, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contenido no encontrado"})
	}

	return c.JSON(http.StatusOK, content)
}
// UpdateContent actualiza contenido existente
// @Summary Update content
// @Description Update premium content (creator only)
// @Tags Content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param request body ContentRequest true "Update data"
// @Success 200 {object} models.PremiumContent
// @Router /content/{id} [put]
func UpdateContent(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || !user.IsCreator {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Solo creadores pueden actualizar contenido"})
	}

	contentID := c.Param("id")
	
	var req ContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos inválidos"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	db := c.Get("db").(*gorm.DB)
	
	// Verificar que el contenido existe y pertenece al creador
	var content models.PremiumContent
	if err := db.First(&content, contentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contenido no encontrado"})
	}

	if content.CreatorID != user.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "No eres el creador de este contenido"})
	}

	// Actualizar campos
	updates := map[string]interface{}{
		"title":        req.Title,
		"description":  req.Description,
		"price":        req.Price,
		"content_type": req.ContentType,
		"is_locked":    req.IsLocked,
	}

	if err := db.Model(&content).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al actualizar contenido"})
	}

	return c.JSON(http.StatusOK, content)
}

// DeleteContent elimina contenido
// @Summary Delete content
// @Description Delete premium content (creator only)
// @Tags Content
// @Param id path int true "Content ID"
// @Success 204
// @Router /content/{id} [delete]
func DeleteContent(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || !user.IsCreator {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Solo creadores pueden eliminar contenido"})
	}

	contentID := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	
	// Verificar que el contenido existe y pertenece al creador
	var content models.PremiumContent
	if err := db.First(&content, contentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contenido no encontrado"})
	}

	if content.CreatorID != user.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "No eres el creador de este contenido"})
	}

	if err := db.Delete(&content).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al eliminar contenido"})
	}

	return c.NoContent(http.StatusNoContent)
}

type SearchContentParams struct {
	Query string `query:"q"`
	Type  string `query:"type"`
	Tag   string `query:"tag"`
	Limit int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Page  int    `query:"page" validate:"omitempty,min=1"`
}

// SearchContent busca contenido con filtros
// @Summary Search content
// @Description Search premium content with filters
// @Tags Content
// @Produce json
// @Param q query string false "Search query"
// @Param type query string false "Content type (video, image, podcast, document)"
// @Param tag query string false "Tag name"
// @Param limit query int false "Results per page (max 100)" default(20)
// @Param page query int false "Page number" default(1)
// @Success 200 {array} models.PremiumContent
// @Router /content/search [get]
func SearchContent(c echo.Context) error {
	var params SearchContentParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Parámetros inválidos"})
	}

	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validación fallida", "details": err.Error()})
	}

	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Page == 0 {
		params.Page = 1
	}

	db := c.Get("db").(*gorm.DB)
	query := db.Model(&models.PremiumContent{}).
		Where("is_locked = ?", false).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit)

	if params.Query != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", 
			"%"+params.Query+"%", "%"+params.Query+"%")
	}

	if params.Type != "" {
		query = query.Where("content_type = ?", params.Type)
	}

	if params.Tag != "" {
		query = query.Joins("JOIN content_tags ON content_tags.content_id = premium_contents.id").
			Joins("JOIN tags ON tags.id = content_tags.tag_id").
			Where("tags.name = ?", params.Tag)
	}

	var results []models.PremiumContent
	if err := query.Find(&results).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al buscar contenido"})
	}

	return c.JSON(http.StatusOK, results)
}
// AddTagsToContent - Asignar tags a contenido
// @Summary Add tags to content
// @Description Assign tags to specific content
// @Tags Content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param tags body []string true "Array of tag names"
// @Success 200 {object} map[string]interface{}
// @Router /content/{id}/tags [post]
func AddTagsToContent(c echo.Context) error {
    // Obtener usuario autenticado
    user, ok := c.Get("user").(*models.User)
    if !ok {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "error": "Usuario no autenticado",
        })
    }

    contentID := c.Param("id")
    db := c.Get("db").(*gorm.DB)

    // 1. Verificar que el contenido existe y pertenece al usuario
    var content models.PremiumContent
    if err := db.First(&content, contentID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]interface{}{
            "error": "Contenido no encontrado",
        })
    }

    if content.CreatorID != user.ID {
        return c.JSON(http.StatusForbidden, map[string]interface{}{
            "error": "No tienes permiso para modificar este contenido",
        })
    }

    // 2. Parsear los tags del body
    var tagNames []string
    if err := c.Bind(&tagNames); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "error": "Formato de tags inválido",
            "details": "Debe ser un array de strings",
        })
    }

    // 3. Validar cantidad de tags (máximo 5 por contenido)
    if len(tagNames) > 5 {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "error": "Demasiadas etiquetas",
            "details": "Máximo 5 etiquetas por contenido",
        })
    }

    // 4. Procesar cada tag
    var assignedTags []models.Tag
    for _, name := range tagNames {
        // Limpiar y validar el nombre del tag
        name = strings.TrimSpace(name)
        if name == "" {
            continue
        }

        // Buscar o crear el tag (insensible a mayúsculas)
        var tag models.Tag
        err := db.Where("LOWER(name) = LOWER(?)", name).First(&tag).Error
        
        if err != nil {
            // Tag no existe, crearlo
            tag = models.Tag{
                Name: name,
                Slug: slug.Make(name),
            }
            if err := db.Create(&tag).Error; err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                    "error": "Error al crear etiqueta",
                    "details": err.Error(),
                })
            }
        }

        // Verificar si la relación ya existe
        var existingRel models.ContentTag
        if err := db.Where("content_id = ? AND tag_id = ?", content.ID, tag.ID).First(&existingRel).Error; err != nil {
            // Crear relación si no existe
            rel := models.ContentTag{
                ContentID: content.ID,
                TagID:     tag.ID,
            }
            if err := db.Create(&rel).Error; err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                    "error": "Error al asignar etiqueta",
                    "details": err.Error(),
                })
            }
        }

        assignedTags = append(assignedTags, tag)
    }

    // 5. Obtener todos los tags del contenido para la respuesta
    var currentTags []models.Tag
    db.Joins("JOIN content_tags ON content_tags.tag_id = tags.id").
        Where("content_tags.content_id = ?", content.ID).
        Find(&currentTags)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Etiquetas asignadas correctamente",
        "content_id": content.ID,
        "assigned_tags": assignedTags,
        "all_tags": currentTags,
    })
}