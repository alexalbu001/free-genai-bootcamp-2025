package activities

import (
	"net/http"
	"strconv"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	activityService *service.ActivityService
}

func NewHandler(activityService *service.ActivityService) *Handler {
	return &Handler{
		activityService: activityService,
	}
}

// RegisterRoutes registers all routes for study activities
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	activities := r.Group("/study_activities")
	{
		activities.GET("/:id", h.Get)
		activities.GET("/:id/study_sessions", h.ListSessions)
	}
}

// Get returns a single study activity
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	activity, err := h.activityService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// ListSessions returns study sessions for an activity
func (h *Handler) ListSessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	sessions, total, err := h.activityService.ListSessions(id, page, perPage)
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
