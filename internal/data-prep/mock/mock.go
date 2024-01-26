package mock

import (
	"fmt"
	"net/http"
	"people-service/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func MockServices(cfg *config.Config) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/age", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
		{
			"count": 3800,
			"name": "Dmitriy",
			"age": 43
		}`)
	})

	router.Get("/gender", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"count": 25459,
			"name": "Dmitriy",
			"gender": "male",
			"probability": 1
		}`)
	})

	router.Get("/nat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"count": 34878,
			"name": "Alisa",
			"country": [
				{
					"country_id": "RU",
					"probability": 0.089
				},
				{
					"country_id": "UA",
					"probability": 0.085
				},
				{
					"country_id": "CN",
					"probability": 0.072
				},
				{
					"country_id": "BA",
					"probability": 0.056
				},
				{
					"country_id": "TH",
					"probability": 0.055
				}
			]
		}`)
	})

	srv := &http.Server{
		Addr:         "localhost:8098",
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("failed to start server")
		}
	}()
}
