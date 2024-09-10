package main

import (
	"database/sql"
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

	store := sqlc.New(db)
	server := api.NewServer(store)
	sserver := api.NewStrictHandler(server, nil)
	mux := http.NewServeMux()
	api.HandlerFromMux(sserver, mux)

	s := &http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8081",
	}

	log.Fatal(s.ListenAndServe())
}

func runMigration(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", "mysql", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
