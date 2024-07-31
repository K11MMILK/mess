package storage

import (
	"database/sql"
	"log"

	"mess/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) SaveMessage(msg models.Message) error {
	query := "INSERT INTO messages (content, status) VALUES ($1, $2)"
	_, err := s.DB.Exec(query, msg.Content, msg.Status)
	return err
}

func (s *Storage) UpdateMessageStatus(content, status string) error {
	log.Printf("Updating message status: content = %s, status = %s", content, status)
	query := "UPDATE messages SET status = $1 WHERE content = $2"
	res, err := s.DB.Exec(query, status, content)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Rows affected: %d", rowsAffected)
	return nil
}

func (s *Storage) GetProcessedMessagesCount() (int, error) {
	query := "SELECT COUNT(*) FROM messages WHERE status = 'processed'"
	var count int
	err := s.DB.QueryRow(query).Scan(&count)
	return count, err
}
