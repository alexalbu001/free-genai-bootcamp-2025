//go:build mage

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "words.db"

// InitDB initializes the SQLite database
func InitDB() error {
	if _, err := os.Stat(dbName); err == nil {
		fmt.Printf("Database %s already exists\n", dbName)
		return nil
	}

	file, err := os.Create(dbName)
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	file.Close()

	fmt.Printf("Created database %s\n", dbName)
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("error finding migration files: %v", err)
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Running migration %s\n", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", file, err)
		}

		statements := strings.Split(string(content), ";")

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("error executing migration %s: %v", file, err)
			}
		}
	}

	return nil
}

// Seed imports seed data into the database
func Seed() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return fmt.Errorf("error finding seed files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("Processing seed file %s\n", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading seed file %s: %v", file, err)
		}

		var seedData struct {
			GroupName string `json:"group_name"`
			Words     []struct {
				French  string `json:"french"`
				English string `json:"english"`
			} `json:"words"`
		}

		if err := json.Unmarshal(content, &seedData); err != nil {
			return fmt.Errorf("error parsing seed file %s: %v", file, err)
		}

		// Create group
		result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", seedData.GroupName)
		if err != nil {
			return fmt.Errorf("error creating group: %v", err)
		}

		groupID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting group ID: %v", err)
		}

		// Insert words and create word_group relationships
		for _, word := range seedData.Words {
			parts, err := json.Marshal(map[string]string{
				"french":  word.French,
				"english": word.English,
			})
			if err != nil {
				return fmt.Errorf("error marshaling word parts: %v", err)
			}

			result, err := db.Exec("INSERT INTO words (parts) VALUES (?)", parts)
			if err != nil {
				return fmt.Errorf("error creating word: %v", err)
			}

			wordID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("error getting word ID: %v", err)
			}

			_, err = db.Exec("INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)", wordID, groupID)
			if err != nil {
				return fmt.Errorf("error creating word_group: %v", err)
			}
		}

		// Update words_count
		_, err = db.Exec("UPDATE groups SET words_count = ? WHERE id = ?", len(seedData.Words), groupID)
		if err != nil {
			return fmt.Errorf("error updating words_count: %v", err)
		}
	}

	// Process study activities
	activityFiles, err := filepath.Glob("db/seeds/study_activities.json")
	if err != nil {
		return fmt.Errorf("error finding activity seed files: %v", err)
	}

	for _, file := range activityFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading activity seed file %s: %v", file, err)
		}

		var activityData struct {
			Activities []struct {
				Name         string `json:"name"`
				URL          string `json:"url"`
				ThumbnailURL string `json:"thumbnail_url"`
				Description  string `json:"description"`
			} `json:"activities"`
		}

		if err := json.Unmarshal(content, &activityData); err != nil {
			return fmt.Errorf("error parsing activity seed file %s: %v", file, err)
		}

		for _, activity := range activityData.Activities {
			_, err := db.Exec(`
				INSERT INTO study_activities (name, url, thumbnail_url, description)
				VALUES (?, ?, ?, ?)
			`, activity.Name, activity.URL, activity.ThumbnailURL, activity.Description)
			if err != nil {
				return fmt.Errorf("error creating activity: %v", err)
			}
		}
	}

	return nil
}

// Reset resets all data in the database
func Reset() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM word_review_items;
		DELETE FROM study_sessions;
		DELETE FROM word_groups;
		DELETE FROM words;
		DELETE FROM groups;
		DELETE FROM study_activities;
	`)

	return err
}
