package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
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

		-- Insert some test data
		INSERT INTO words (parts) VALUES 
			('{"french":"bonjour","english":"hello"}'),
			('{"french":"au revoir","english":"goodbye"}'),
			('{"french":"merci","english":"thank you"}');

		-- Insert some groups
		INSERT INTO groups (name, words_count) VALUES 
			('Basics', 2),
			('Greetings', 1);

		-- Link words to groups
		INSERT INTO word_groups (word_id, group_id) VALUES 
			(1, 1),
			(2, 1),
			(1, 2);

		-- Insert study activities
		INSERT INTO study_activities (name, url, thumbnail_url, description) VALUES
			('Flashcards', '/activities/flashcards', '/thumbnails/flashcards.jpg', 'Practice with flashcards'),
			('Quiz', '/activities/quiz', '/thumbnails/quiz.jpg', 'Test your knowledge');
	`

	_, err = db.Exec(migrationSQL)
	if err != nil {
		return err
	}

	return db.Ping()
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// SetDB sets the database instance
func SetDB(newDB *sql.DB) {
	db = newDB
}
