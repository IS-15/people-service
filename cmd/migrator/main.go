package main

import (
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// IS: create database before start
func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("start migrations")

	m, err := migrate.New(
		"file://"+"./migrations",
		"postgres://postgres:postgres@localhost:5432/person-service?sslmode=disable")
	if err != nil {
		log.Error("error occured while migration init: ", err)
	}

	log.Info("start migration up")

	if err := m.Up(); err != nil {
		log.Error("error occured while migration up: ", err)
	}

	log.Info("migration finished")
}
