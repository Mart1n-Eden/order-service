package storage

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddOrder(id string, content json.RawMessage) error {
	if _, err := s.db.Exec("INSERT INTO orders (content,uid) VALUES ($1,$2)", content, id); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetOrder(id string) (content json.RawMessage, err error) {
	row := s.db.QueryRow("SELECT content FROM orders WHERE uid = $1", id)
	row.Scan(&content)

	// TODO: handling error

	return content, nil
}

func (s *Storage) FillCache() (ids []string, contents []json.RawMessage, err error) {
	rows, err := s.db.Query("SELECT uid, content FROM orders")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   string
			cont json.RawMessage
		)

		if err := rows.Scan(&id, &cont); err != nil {
			return nil, nil, err
		}

		ids = append(ids, id)
		contents = append(contents, cont)
	}

	return ids, contents, nil
}
