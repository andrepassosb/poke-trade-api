package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/andrepassosb/poke-trade-api/database"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
	cardUpdateQueue chan CardUpdate // Adicionando o canal
}

func Main() {
	url := database.GetURL()

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

  	models := database.NewModels(db)
	app := &application{
		port:      8080,
		jwtSecret: os.Getenv("JWT_SECRET"),
		models:    models,
		cardUpdateQueue: make(chan CardUpdate, 100), // Buffer de 100 updates
	}

	go app.startCardUpdateWorker()

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
