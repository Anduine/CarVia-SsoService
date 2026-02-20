package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(host, dbName, user, password string) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName))
	if err != nil {
		slog.Error("Problem with db connection:", "Error", err)
		db.Close()
		return nil
	}

	if err := db.Ping(); err != nil {
		slog.Error("Problem with ping db:", "Error", err)
	}

	counts := 0
	for {
		err := db.Ping()
		if err == nil {
			slog.Info("Successfully connected to Postgres")
			break
		}

		slog.Warn("Postgres not ready...", "count", counts, "err", err)
		counts++

		if counts > 5 {
			slog.Error("Could not connect to DB after many retries")
			os.Exit(1)
		}

		slog.Info("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	return db
}
