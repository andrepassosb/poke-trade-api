package database

import (
	"log"
	"os"
)

func GetURL() string {
	environment := "local"

	if env := os.Getenv("DATABASE_ENV"); env != "" {
		environment = env
	}

	var dbURL string
	switch environment {
	case "local":
		dbURL = "file:local.db"
	case "TURSO_DATABASE":
		url := os.Getenv("TURSO_DATABASE_URL")
		token := os.Getenv("TURSO_AUTH_TOKEN")
		if url == "" || token == "" {
			log.Fatal("TURSO_DATABASE_URL and TURSO_AUTH_TOKEN must be set for remote environment")
		}
		dbURL = url + "?auth_token=" + token
	default:
		log.Fatal("Invalid environment. Use 'local' or 'TURSO_PRODUCTION'.")
	}

	return dbURL

}