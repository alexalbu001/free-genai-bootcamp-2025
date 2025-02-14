package service

import (
	"database/sql"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
)

type DashboardService struct {
	db *sql.DB
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		db: storage.GetDB(),
	}
}

type LastStudySession struct {
	ID              int64  `json:"id"`
	GroupID         int64  `json:"group_id"`
	StudyActivityID int64  `json:"study_activity_id"`
	GroupName       string `json:"group_name"`
}

type StudyProgress struct {
	TotalWordsStudied   int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

type QuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

// GetLastStudySession returns the most recent study session
func (s *DashboardService) GetLastStudySession() (*LastStudySession, error) {
	var session LastStudySession
	err := s.db.QueryRow(`
		SELECT s.id, s.group_id, s.study_activity_id, g.name
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`).Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.GroupName)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetStudyProgress returns the overall study progress
func (s *DashboardService) GetStudyProgress() (*StudyProgress, error) {
	var progress StudyProgress

	// Get total available words
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM words
	`).Scan(&progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	// Get total words studied (unique words reviewed)
	err = s.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&progress.TotalWordsStudied)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

// GetQuickStats returns quick statistics about the user's learning
func (s *DashboardService) GetQuickStats() (*QuickStats, error) {
	var stats QuickStats

	// Get success rate
	err := s.db.QueryRow(`
		SELECT COALESCE(
			(SELECT CAST(SUM(CASE WHEN correct THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) * 100
			FROM word_review_items), 0)
	`).Scan(&stats.SuccessRate)
	if err != nil {
		return nil, err
	}

	// Get total study sessions
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM study_sessions
	`).Scan(&stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	// Get total active groups (groups with at least one study session)
	err = s.db.QueryRow(`
		SELECT COUNT(DISTINCT group_id) 
		FROM study_sessions
	`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	// Get study streak (consecutive days with study sessions)
	err = s.db.QueryRow(`
		WITH RECURSIVE dates AS (
			SELECT date(created_at) as study_date
			FROM study_sessions
			GROUP BY date(created_at)
			ORDER BY study_date DESC
			LIMIT 1
		),
		streak AS (
			SELECT study_date, 1 as days
			FROM dates
			UNION ALL
			SELECT date(study_date, '-1 day'), days + 1
			FROM streak
			WHERE EXISTS (
				SELECT 1
				FROM study_sessions
				WHERE date(created_at) = date(study_date, '-1 day')
			)
		)
		SELECT COUNT(*) FROM streak
	`).Scan(&stats.StudyStreakDays)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
