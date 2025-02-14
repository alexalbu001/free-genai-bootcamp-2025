package models

import (
	"encoding/json"
	"time"
)

// Group represents a collection of words
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WordsCount int    `json:"words_count"`
}

// StudyActivity represents a learning activity type
type StudyActivity struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

// StudySession represents a learning session
type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// WordReviewItem represents a single word review in a study session
type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// Word represents a word in the group
type Word struct {
	ID    int64           `json:"id"`
	Parts json.RawMessage `json:"parts"`
}

// ScanWord is a helper type for scanning JSON from database
type ScanWord struct {
	ID    int64  `json:"id"`
	Parts []byte `json:"parts"`
}

// ToWord converts ScanWord to Word
func (sw *ScanWord) ToWord() Word {
	return Word{
		ID:    sw.ID,
		Parts: json.RawMessage(sw.Parts),
	}
}

// WordParts represents the parts of a word
type WordParts struct {
	French  string `json:"french"`
	English string `json:"english"`
}
