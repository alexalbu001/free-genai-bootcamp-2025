package service

import (
	"database/sql"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/models"
	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
)

type GroupService struct {
	db *sql.DB
}

func NewGroupService() *GroupService {
	return &GroupService{
		db: storage.GetDB(),
	}
}

// List returns a paginated list of groups
func (s *GroupService) List(page, perPage int) ([]models.Group, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT id, name, words_count 
		FROM groups 
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.Name, &group.WordsCount)
		if err != nil {
			return nil, 0, err
		}
		groups = append(groups, group)
	}

	return groups, total, nil
}

// Get returns a single group by ID
func (s *GroupService) Get(id int64) (*models.Group, error) {
	var group models.Group
	err := s.db.QueryRow(`
		SELECT id, name, words_count 
		FROM groups 
		WHERE id = ?
	`, id).Scan(&group.ID, &group.Name, &group.WordsCount)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// ListWords returns words in a group
func (s *GroupService) ListWords(groupID int64, page, perPage int) ([]models.Word, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM word_groups 
		WHERE group_id = ?
	`, groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT w.id, json(w.parts) as parts
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		LIMIT ? OFFSET ?
	`, groupID, perPage, offset)
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

// ListStudySessions returns study sessions for a group
func (s *GroupService) ListStudySessions(groupID int64, page, perPage int) ([]models.StudySession, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE group_id = ?
	`, groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, groupID, perPage, offset)
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
