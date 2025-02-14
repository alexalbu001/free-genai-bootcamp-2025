package service

import (
	"database/sql"
	"time"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/models"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
)

type SessionService struct {
	db *sql.DB
}

func NewSessionService() *SessionService {
	return &SessionService{
		db: storage.GetDB(),
	}
}

// List returns a paginated list of study sessions
func (s *SessionService) List(page, perPage int) ([]SessionResponse, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			strftime('%Y-%m-%d %H:%M:%S', ss.created_at) as start_time,
			strftime('%Y-%m-%d %H:%M:%S', COALESCE(MAX(wri.created_at), ss.created_at)) as end_time,
			COUNT(wri.id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []SessionResponse
	for rows.Next() {
		var session SessionResponse
		err := rows.Scan(
			&session.ID,
			&session.ActivityName,
			&session.GroupName,
			&session.StartTime,
			&session.EndTime,
			&session.ReviewItemsCount,
		)
		if err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}

// Create creates a new study session
func (s *SessionService) Create(groupID, studyActivityID int64) (*SessionResponse, error) {
	result, err := s.db.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, ?)
	`, groupID, studyActivityID, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.Get(id)
}

// Add this struct for the session response
type SessionResponse struct {
	ID               int64  `json:"id"`
	ActivityName     string `json:"activity_name"`
	GroupName        string `json:"group_name"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	ReviewItemsCount int    `json:"review_items_count"`
}

// Update the Get method to include these fields
func (s *SessionService) Get(id int64) (*SessionResponse, error) {
	var session SessionResponse
	err := s.db.QueryRow(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			strftime('%Y-%m-%d %H:%M:%S', ss.created_at) as start_time,
			strftime('%Y-%m-%d %H:%M:%S', COALESCE(MAX(wri.created_at), ss.created_at)) as end_time,
			COUNT(wri.id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.id = ?
		GROUP BY ss.id
	`, id).Scan(
		&session.ID,
		&session.ActivityName,
		&session.GroupName,
		&session.StartTime,
		&session.EndTime,
		&session.ReviewItemsCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// ReviewWord records a word review in a study session
func (s *SessionService) ReviewWord(sessionID, wordID int64, correct bool) (*models.WordReviewItem, error) {
	result, err := s.db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`, wordID, sessionID, correct, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var review models.WordReviewItem
	err = s.db.QueryRow(`
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE id = ?
	`, id).Scan(&review.ID, &review.WordID, &review.StudySessionID, &review.Correct, &review.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// ListWords returns words reviewed in a study session
func (s *SessionService) ListWords(sessionID int64, page, perPage int) ([]models.Word, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(DISTINCT w.id)
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
	`, sessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT DISTINCT w.id, json(w.parts) as parts
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		LIMIT ? OFFSET ?
	`, sessionID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var scanWord models.ScanWord
		err := rows.Scan(&scanWord.ID, &scanWord.Parts)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, scanWord.ToWord())
	}

	return words, total, nil
}
