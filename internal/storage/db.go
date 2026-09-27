package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// var db *sql.DB

func InitDB() (*sql.DB, error) {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// sql.Open - функция подключения к DB
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort))
	if err != nil {
		return nil, err
	}

	// db.Ping() - функция проверки успешности подключения к DB
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
