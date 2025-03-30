package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}
	direction := os.Args[1]
	environment := "local" // valor padrão


	if env := os.Getenv("DATABASE_ENV"); env != "" {
		environment = env
	}


	var dbURL string
	switch environment {
	case "local":
		dbURL = "file:local.db" // URL do banco de dados local
	case "TURSO_PRODUCTION":
		url := os.Getenv("TURSO_DATABASE_URL")
		token := os.Getenv("TURSO_AUTH_TOKEN")
		if url == "" || token == "" {
			log.Fatal("TURSO_DATABASE_URL and TURSO_AUTH_TOKEN must be set for remote environment")
		}
		dbURL = url + "?auth_token=" + token
	default:
		log.Fatal("Invalid environment. Use 'local' or 'TURSO_PRODUCTION'.")
	}


	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
}
