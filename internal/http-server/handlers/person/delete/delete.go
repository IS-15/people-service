package delete

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	resp "people-service/internal/lib/api/response"
	"people-service/internal/lib/routing"
)

type Response struct {
	resp.Response
}

type PersonDeleter interface {
	DeletePerson(id int) error
}

func New(log *slog.Logger, personDeleter PersonDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, routing.PersonIdParam)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Info("error while parsing person id", slog.String("id", idParam))
			render.JSON(w, r, resp.Error("error while deleting person"))
			return
		}

		err = personDeleter.DeletePerson(id)
		if err != nil {
			log.Info("error while deleting person", slog.Int("id", id))
			render.JSON(w, r, resp.Error("error while deleting person"))
			return
		}

		log.Info("person deleted", slog.Int("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
