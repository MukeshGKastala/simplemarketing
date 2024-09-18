package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MukeshGKastala/marketing/api"
	sqlc "github.com/MukeshGKastala/marketing/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := runMigration(db); err != nil {
		log.Fatal(err)
	}

	store := sqlc.NewStore(db)
	server := api.NewStrictHandler(api.NewServer(store), nil)
	mux := http.NewServeMux()
	api.HandlerFromMuxWithBaseURL(server, mux, "/marketing")

	// TODO: Configure timeouts
	s := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")),
	}

	// TODO: Implement signal handling and graceful shutdown
	log.Fatal(s.ListenAndServe())
}

// TODO: Refactor to accept DSN parts (host, user, etc.) to support connecting to multiple databases and ease testing.
func runMigration(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migration", "mysql", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
