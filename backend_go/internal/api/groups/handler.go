package groups

import (
	"net/http"
	"strconv"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	groupService *service.GroupService
}

func NewHandler(groupService *service.GroupService) *Handler {
	return &Handler{
		groupService: groupService,
	}
}

// RegisterRoutes registers all routes for groups
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	groups := r.Group("/groups")
	{
		groups.GET("", h.List)
		groups.GET("/:id", h.Get)
		groups.GET("/:id/words", h.ListWords)
		groups.GET("/:id/study_sessions", h.ListStudySessions)
	}
}

// List returns a paginated list of groups
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	groups, total, err := h.groupService.List(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + perPage - 1) / perPage,
			"total_items":    total,
			"items_per_page": perPage,
		},
	})
}

// Get returns a single group by ID
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.groupService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// ListWords returns words in a group
func (h *Handler) ListWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	words, total, err := h.groupService.ListWords(id, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": words,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + perPage - 1) / perPage,
			"total_items":    total,
			"items_per_page": perPage,
		},
	})
}

// ListStudySessions returns study sessions for a group
func (h *Handler) ListStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	sessions, total, err := h.groupService.ListStudySessions(id, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sessions,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + perPage - 1) / perPage,
			"total_items":    total,
			"items_per_page": perPage,
		},
	})
}
