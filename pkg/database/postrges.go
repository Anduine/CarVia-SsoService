package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(dbConn string) (*sql.DB) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal("Problem with db connection:", err)
		db.Close()
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Problem with ping db:", err)
	}

	log.Println("DB Connection successfully")
	return db
}