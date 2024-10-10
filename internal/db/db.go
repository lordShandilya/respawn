package db

import (
	"backend/internal/config"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	connectionStr := config.LoadConfig().DatabaseURL

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("Database connection failed!!")
		return errors.New("Database connection error.")
	}

	DB = db

	return nil

}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
