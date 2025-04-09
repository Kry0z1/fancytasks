package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgresql", 5432, "postgresql", os.Getenv("strong_postgres_pass"), "default",
	)

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Couldn't open database: %s", err.Error())
	}
}

func GetDB() *sql.DB {
	return db
}
