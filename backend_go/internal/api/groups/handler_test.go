package groups

import (
	"database/sql"
	"encoding/json"
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

	groupService := service.NewGroupService()
	handler := NewHandler(groupService)

	r := gin.New()
	api := r.Group("/api")
	handler.RegisterRoutes(api)

	return r, db
}

func TestListGroups(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO groups (name, words_count) VALUES 
        ('Group 1', 5),
        ('Group 2', 3)
    `)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/groups", nil)
	w := testutil.ExecuteRequest(r, req)
	testutil.LogResponse(t, w)

	testutil.CheckResponseCode(t, http.StatusOK, w.Code)

	var response struct {
		Items []struct {
			ID         int64  `json:"id"`
			Name       string `json:"name"`
			WordsCount int    `json:"words_count"`
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
		t.Errorf("Expected 2 groups, got %d", len(response.Items))
	}
}

func TestListGroupWords(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO groups (name) VALUES ('Test Group');
        INSERT INTO words (parts) VALUES ('{"french":"test1","english":"test1"}');
        INSERT INTO word_groups (word_id, group_id) VALUES (1, 1);
    `)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/groups/1/words", nil)
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
