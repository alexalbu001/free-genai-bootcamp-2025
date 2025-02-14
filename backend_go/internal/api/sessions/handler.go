package sessions

import (
	"net/http"
	"strconv"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	sessionService *service.SessionService
}

func NewHandler(sessionService *service.SessionService) *Handler {
	return &Handler{
		sessionService: sessionService,
	}
}

// RegisterRoutes registers all routes for sessions
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	sessions := r.Group("/study_sessions")
	{
		sessions.GET("", h.List)
		sessions.GET("/:id", h.Get)
		sessions.GET("/:id/words", h.ListWords)
		sessions.POST("", h.Create)
		sessions.POST("/:id/word/:word_id/review", h.ReviewWord)
	}
}

// List returns a paginated list of study sessions
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	sessions, total, err := h.sessionService.List(page, perPage)
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

// Create creates a new study session
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		GroupID         int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionService.Create(req.GroupID, req.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// ReviewWord records a word review in a study session
func (h *Handler) ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word id"})
		return
	}

	var req struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.sessionService.ReviewWord(sessionID, wordID, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// Get returns a single study session
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	session, err := h.sessionService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// ListWords returns words reviewed in a study session
func (h *Handler) ListWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	words, total, err := h.sessionService.ListWords(id, page, perPage)
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
