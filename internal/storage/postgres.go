package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"order-service/internal/config"
)

func NewPostgres(cfg config.DBConfig) (db *sql.DB, err error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	if db, err = sql.Open("postgres", url); err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
