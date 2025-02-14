package words

import (
	"net/http"
	"strconv"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	wordService *service.WordService
}

func NewHandler(wordService *service.WordService) *Handler {
	return &Handler{
		wordService: wordService,
	}
}

// RegisterRoutes registers all routes for words
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	words := r.Group("/words")
	{
		words.GET("", h.List)
		words.GET("/:id", h.Get)
	}
}

// List returns a paginated list of words
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	words, total, err := h.wordService.List(page, perPage)
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

// Get returns a single word by ID
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	word, err := h.wordService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}
