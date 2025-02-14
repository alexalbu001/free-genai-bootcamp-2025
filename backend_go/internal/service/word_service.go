package service

import (
	"database/sql"
	"encoding/json"

	"github.com/free-genai-bootcamp-2025/lang-portal/backend_go/internal/storage"
)

type WordService struct {
	db *sql.DB
}

func NewWordService() *WordService {
	return &WordService{
		db: storage.GetDB(),
	}
}

// Add this struct for word response
type WordResponse struct {
	ID           int64           `json:"id"`
	Parts        json.RawMessage `json:"parts"`
	CorrectCount int             `json:"correct_count"`
	WrongCount   int             `json:"wrong_count"`
}

// Update the List method to include counts
func (s *WordService) List(page, perPage int) ([]WordResponse, int, error) {
	offset := (page - 1) * perPage

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT 
			w.id, 
			json(w.parts) as parts,
			COALESCE(SUM(CASE WHEN wri.correct THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN NOT wri.correct THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []WordResponse
	for rows.Next() {
		var scanWord struct {
			ID           int64  `json:"id"`
			Parts        []byte `json:"parts"`
			CorrectCount int    `json:"correct_count"`
			WrongCount   int    `json:"wrong_count"`
		}
		err := rows.Scan(&scanWord.ID, &scanWord.Parts, &scanWord.CorrectCount, &scanWord.WrongCount)
		if err != nil {
			return nil, 0, err
		}
		words = append(words, WordResponse{
			ID:           scanWord.ID,
			Parts:        json.RawMessage(scanWord.Parts),
			CorrectCount: scanWord.CorrectCount,
			WrongCount:   scanWord.WrongCount,
		})
	}

	return words, total, nil
}

// Update the Get method to include counts
func (s *WordService) Get(id int64) (*WordResponse, error) {
	var scanWord struct {
		ID           int64  `json:"id"`
		Parts        []byte `json:"parts"`
		CorrectCount int    `json:"correct_count"`
		WrongCount   int    `json:"wrong_count"`
	}
	err := s.db.QueryRow(`
		SELECT 
			w.id, 
			json(w.parts) as parts,
			COALESCE(SUM(CASE WHEN wri.correct THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN NOT wri.correct THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`, id).Scan(&scanWord.ID, &scanWord.Parts, &scanWord.CorrectCount, &scanWord.WrongCount)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &WordResponse{
		ID:           scanWord.ID,
		Parts:        json.RawMessage(scanWord.Parts),
		CorrectCount: scanWord.CorrectCount,
		WrongCount:   scanWord.WrongCount,
	}, nil
}
