package dashboard

import (
	"net/http"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dashboardService *service.DashboardService
}

func NewHandler(dashboardService *service.DashboardService) *Handler {
	return &Handler{
		dashboardService: dashboardService,
	}
}

// RegisterRoutes registers all dashboard routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/last_study_session", h.LastStudySession)
		dashboard.GET("/study_progress", h.StudyProgress)
		dashboard.GET("/quick_stats", h.QuickStats)
	}
}

// LastStudySession returns the most recent study session
func (h *Handler) LastStudySession(c *gin.Context) {
	session, err := h.dashboardService.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// StudyProgress returns the overall study progress
func (h *Handler) StudyProgress(c *gin.Context) {
	progress, err := h.dashboardService.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// QuickStats returns quick statistics about the user's learning
func (h *Handler) QuickStats(c *gin.Context) {
	stats, err := h.dashboardService.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
