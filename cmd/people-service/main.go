package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"people-service/config"
	"people-service/internal/data-prep/age"
	"people-service/internal/data-prep/gender"
	"people-service/internal/data-prep/nationality"
	"people-service/internal/http-server/handlers/person/delete"
	"people-service/internal/http-server/handlers/person/get"
	"people-service/internal/http-server/handlers/person/save"
	"people-service/internal/http-server/handlers/person/update"
	mwLogger "people-service/internal/http-server/middleware/logger"
	"people-service/internal/lib/routing"
	"people-service/internal/storage"
	"people-service/internal/storage/pg"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// // IS: Use for easy init env variables. Not for production use. Only for study case.
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {

	cfg := config.MustLoad()

	pgConfig := storage.PostgresConfig{
		Host:     cfg.Storage.Host,
		Port:     cfg.Storage.Port,
		DBName:   cfg.Storage.DBName,
		User:     cfg.Storage.User,
		Password: cfg.Storage.Password,
	}

	log := setupLogger(cfg.Env)

	log.Debug("init database")
	storage, err := pg.New(log, pgConfig)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}
	log.Debug("db initialized")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.CtxTimeout)*time.Second)
	defer cancel()

	log.Debug("initializing data preparation services")
	ageService := age.New(log, cfg.AgeServiceUrl) // mock: "http://localhost:8098/age"
	log.Debug("age service initialized")
	genderService := gender.New(log, cfg.GenderServiceUrl) // mock: "http://localhost:8098/gender")
	log.Debug("gender service initialized")
	nationalityService := nationality.New(log, cfg.NationalityServiceUrl) // mock: "http://localhost:8098/nat"
	log.Debug("nationality service initialized")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/person", save.New(log, storage, ageService, genderService, nationalityService))
	router.Get("/person", get.New(log, storage))

	router.Route(fmt.Sprintf("/person/{%s}", routing.PersonIdParam), func(r chi.Router) {
		r.Delete("/", delete.New(log, storage))
		r.Put("/", update.New(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// IS: if have no access to real services, simple mock for them.
	// TODO: better use testing and mocking frameworks
	//mock.MockServices(cfg)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)
		return
	}

	storage.Close()

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
