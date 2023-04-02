package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"))

	db, errOpen := sql.Open("postgres", dsn)
	if errOpen != nil {
		return nil, errOpen
	}

	if errPing := db.Ping(); errPing != nil {
		return nil, errPing
	}

	return db, nil
}
