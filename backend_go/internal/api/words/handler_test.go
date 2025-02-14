package words

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/service"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/testutil"
	"github.com/gin-gonic/gin"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	db := testutil.SetupTestDB(t)
	testutil.SetTestDB(db)

	wordService := service.NewWordService()
	handler := NewHandler(wordService)

	r := gin.New()
	api := r.Group("/api")
	handler.RegisterRoutes(api)

	return r, db
}

func TestListWords(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO words (parts) VALUES
		('{"french":"bonjour","english":"hello"}'),
		('{"french":"au revoir","english":"goodbye"}')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Make request
	req := httptest.NewRequest("GET", "/api/words", nil)
	w := testutil.ExecuteRequest(r, req)

	// Check response
	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var response struct {
		Items []struct {
			ID           int64           `json:"id"`
			Parts        json.RawMessage `json:"parts"`
			CorrectCount int             `json:"correct_count"`
			WrongCount   int             `json:"wrong_count"`
		} `json:"items"`
		Pagination struct {
			CurrentPage  int `json:"current_page"`
			TotalPages   int `json:"total_pages"`
			TotalItems   int `json:"total_items"`
			ItemsPerPage int `json:"items_per_page"`
		} `json:"pagination"`
	}

	testutil.ParseResponse(t, w, &response)

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 words, got %d", len(response.Items))
	}

	// Only check if we have items
	if len(response.Items) > 0 {
		var parts struct {
			French  string `json:"french"`
			English string `json:"english"`
		}
		if err := json.Unmarshal(response.Items[0].Parts, &parts); err != nil {
			t.Fatalf("Failed to parse word parts: %v", err)
		}

		if parts.French != "bonjour" {
			t.Errorf("Expected first word to be 'bonjour', got '%s'", parts.French)
		}
	}
}

func TestGetWord(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	result, err := db.Exec(`
		INSERT INTO words (parts) VALUES ('{"french":"bonjour","english":"hello"}')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	id, _ := result.LastInsertId()

	// Test successful get
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/words/%d", id), nil)
	w := testutil.ExecuteRequest(r, req)

	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var word service.WordResponse
	testutil.ParseResponse(t, w, &word)

	var parts struct {
		French  string `json:"french"`
		English string `json:"english"`
	}
	if err := json.Unmarshal(word.Parts, &parts); err != nil {
		t.Fatalf("Failed to parse word parts: %v", err)
	}

	if parts.French != "bonjour" {
		t.Errorf("Expected word to be 'bonjour', got '%s'", parts.French)
	}

	// Test not found
	req = httptest.NewRequest("GET", "/api/words/999", nil)
	w = testutil.ExecuteRequest(r, req)

	testutil.CheckResponseCode(t, http.StatusNotFound, w.Code)
}
