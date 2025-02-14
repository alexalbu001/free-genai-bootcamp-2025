package testutil

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates a test database and returns a connection
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create temporary database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	migrationSQL := `
		-- Create words table
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			parts TEXT NOT NULL
		);

		-- Create groups table
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			words_count INTEGER DEFAULT 0
		);

		-- Create word_groups table
		CREATE TABLE IF NOT EXISTS word_groups (
			word_id INTEGER,
			group_id INTEGER,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			PRIMARY KEY (word_id, group_id)
		);

		-- Create study_activities table
		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			thumbnail_url TEXT,
			description TEXT
		);

		-- Create study_sessions table
		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER,
			study_activity_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
		);

		-- Create word_review_items table
		CREATE TABLE IF NOT EXISTS word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER,
			study_session_id INTEGER,
			correct BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
		);
	`

	_, err = db.Exec(migrationSQL)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// ExecuteRequest performs a test HTTP request and returns the response
func ExecuteRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// CheckResponseCode verifies the HTTP status code
func CheckResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

// LogResponse logs the response status and body
func LogResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
}

// ParseResponse parses the JSON response into the given struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	LogResponse(t, w)
	err := json.Unmarshal(w.Body.Bytes(), v)
	if err != nil {
		t.Fatalf("Failed to parse response: %v\nResponse body: %s", err, w.Body.String())
	}
}

// Add this function
func SetTestDB(db *sql.DB) {
	storage.SetDB(db)
}

func init() {
	// Disable Gin debug output in tests
	gin.SetMode(gin.ReleaseMode)
}
