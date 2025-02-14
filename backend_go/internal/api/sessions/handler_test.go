package sessions

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

	sessionService := service.NewSessionService()
	handler := NewHandler(sessionService)

	r := gin.New()
	api := r.Group("/api")
	handler.RegisterRoutes(api)

	return r, db
}

func TestCreateSession(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert prerequisite data
	_, err := db.Exec(`
		INSERT INTO groups (name) VALUES ('Test Group');
		INSERT INTO study_activities (name, url, thumbnail_url, description) 
		VALUES ('Test Activity', 'http://test.com', 'http://test.com/thumb.jpg', 'Test Description');
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Create request body
	body := map[string]interface{}{
		"group_id":          1,
		"study_activity_id": 1,
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req := httptest.NewRequest("POST", "/api/study_sessions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := testutil.ExecuteRequest(r, req)

	// Check response
	testutil.CheckResponseCode(t, http.StatusCreated, w.Code)

	var response service.SessionResponse
	testutil.ParseResponse(t, w, &response)

	if response.GroupName != "Test Group" {
		t.Errorf("Expected group name 'Test Group', got '%s'", response.GroupName)
	}
}

func TestReviewWord(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert prerequisite data
	_, err := db.Exec(`
		INSERT INTO words (parts) VALUES ('{"french":"test","english":"test"}');
		INSERT INTO groups (name) VALUES ('Test Group');
		INSERT INTO study_activities (name, url, thumbnail_url, description) 
		VALUES ('Test Activity', 'http://test.com', 'http://test.com/thumb.jpg', 'Test Description');
		INSERT INTO study_sessions (group_id, study_activity_id, created_at) 
		VALUES (1, 1, CURRENT_TIMESTAMP);
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Create request body
	body := map[string]interface{}{
		"correct": true,
	}
	jsonBody, _ := json.Marshal(body)

	// Make request
	req := httptest.NewRequest("POST", "/api/study_sessions/1/word/1/review", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := testutil.ExecuteRequest(r, req)

	// Check response
	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var response struct {
		ID             int64     `json:"id"`
		WordID         int64     `json:"word_id"`
		StudySessionID int64     `json:"study_session_id"`
		Correct        bool      `json:"correct"`
		CreatedAt      time.Time `json:"created_at"`
	}
	testutil.ParseResponse(t, w, &response)

	if !response.Correct {
		t.Error("Expected review to be correct")
	}
}

func TestListSessionWords(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO words (parts) VALUES ('{"french":"test","english":"test"}');
		INSERT INTO groups (name) VALUES ('Test Group');
		INSERT INTO study_activities (name, url, thumbnail_url, description) 
		VALUES ('Test Activity', 'http://test.com', 'http://test.com/thumb.jpg', 'Test Description');
		INSERT INTO study_sessions (group_id, study_activity_id, created_at) 
		VALUES (1, 1, datetime('now'));
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		VALUES (1, 1, true, datetime('now'));
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/study_sessions/1/words", nil)
	w := testutil.ExecuteRequest(r, req)

	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var response struct {
		Items []struct {
			ID    int64           `json:"id"`
			Parts json.RawMessage `json:"parts"`
		} `json:"items"`
	}

	testutil.ParseResponse(t, w, &response)

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 word, got %d", len(response.Items))
	}
}
