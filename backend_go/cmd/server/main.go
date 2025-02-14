package main

import (
	"log"
	"net/http"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/api/activities"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/api/dashboard"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/api/groups"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/api/sessions"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/api/words"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	err := storage.InitDB("words.db") // This should create and initialize the database
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize services
	wordService := service.NewWordService()
	groupService := service.NewGroupService()
	sessionService := service.NewSessionService()
	dashboardService := service.NewDashboardService()
	activityService := service.NewActivityService()

	// Initialize handlers
	wordHandler := words.NewHandler(wordService)
	groupHandler := groups.NewHandler(groupService)
	sessionHandler := sessions.NewHandler(sessionService)
	dashboardHandler := dashboard.NewHandler(dashboardService)
	activityHandler := activities.NewHandler(activityService)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// Register routes for each handler
		wordHandler.RegisterRoutes(api)
		groupHandler.RegisterRoutes(api)
		sessionHandler.RegisterRoutes(api)
		dashboardHandler.RegisterRoutes(api)
		activityHandler.RegisterRoutes(api)

		// Reset endpoints
		api.POST("/reset_history", func(c *gin.Context) {
			db := storage.GetDB()
			_, err := db.Exec(`
				DELETE FROM word_review_items;
				DELETE FROM study_sessions;
			`)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "All study history has been reset",
				"success": true,
			})
		})

		api.POST("/full_reset", func(c *gin.Context) {
			db := storage.GetDB()
			_, err := db.Exec(`
				DELETE FROM word_review_items;
				DELETE FROM study_sessions;
				DELETE FROM word_groups;
				DELETE FROM words;
				DELETE FROM groups;
				DELETE FROM study_activities;
			`)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "System reset complete",
				"success": true,
			})
		})
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
