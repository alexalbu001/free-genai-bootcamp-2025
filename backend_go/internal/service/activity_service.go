package service

import (
	"database/sql"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/models"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
)

type ActivityService struct {
	db *sql.DB
}

func NewActivityService() *ActivityService {
	return &ActivityService{
		db: storage.GetDB(),
	}
}

// Get returns a single study activity by ID
func (s *ActivityService) Get(id int64) (*models.StudyActivity, error) {
	var activity models.StudyActivity
	err := s.db.QueryRow(`
		SELECT id, name, url, thumbnail_url, description
		FROM study_activities
		WHERE id = ?
	`, id).Scan(&activity.ID, &activity.Name, &activity.URL, &activity.ThumbnailURL, &activity.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &activity, nil
}

// ListSessions returns study sessions for an activity
func (s *ActivityService) ListSessions(activityID int64, page, perPage int) ([]models.StudySession, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE study_activity_id = ?
	`, activityID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, activityID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}
