package dashboard

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/testutil"
	"github.com/gin-gonic/gin"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	testutil.SetTestDB(db)

	dashboardService := service.NewDashboardService()
	handler := NewHandler(dashboardService)

	r := gin.New()
	api := r.Group("/api")
	handler.RegisterRoutes(api)

	return r, db
}

func TestGetQuickStats(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO groups (name) VALUES ('Test Group');
		INSERT INTO study_activities (name, url) VALUES ('Test Activity', 'http://test.com');
		INSERT INTO study_sessions (group_id, study_activity_id, created_at) VALUES 
			(1, 1, ?),
			(1, 1, ?);
		INSERT INTO words (parts) VALUES ('{"french":"test","english":"test"}');
		INSERT INTO word_review_items (word_id, study_session_id, correct) VALUES 
			(1, 1, true),
			(1, 2, false);
	`, time.Now(), time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/dashboard/quick_stats", nil)
	w := testutil.ExecuteRequest(r, req)

	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var response struct {
		SuccessRate        float64 `json:"success_rate"`
		TotalStudySessions int     `json:"total_study_sessions"`
		TotalActiveGroups  int     `json:"total_active_groups"`
		StudyStreakDays    int     `json:"study_streak_days"`
	}

	testutil.ParseResponse(t, w, &response)

	if response.SuccessRate != 50.0 {
		t.Errorf("Expected 50%% success rate, got %.2f%%", response.SuccessRate)
	}
	if response.TotalStudySessions != 2 {
		t.Errorf("Expected 2 study sessions, got %d", response.TotalStudySessions)
	}
}
