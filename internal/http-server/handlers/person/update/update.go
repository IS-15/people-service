package update

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"people-service/internal/domain/models"
	resp "people-service/internal/lib/api/response"
	"people-service/internal/lib/routing"
)

type Request struct {
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type Response struct {
	resp.Response
}

type PersonUpdater interface {
	UpdatePerson(id int, person models.Person) error
}

func New(log *slog.Logger, personUpdater PersonUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.update.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idParam := chi.URLParam(r, routing.PersonIdParam)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Info("error while parsing person id", slog.String("id", idParam))
			render.JSON(w, r, resp.Error("error while updating person"))
			return
		}

		var req Request

		err = render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", err)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		person := models.Person{
			Name:        req.Name,
			Surname:     req.Surname,
			Patronymic:  req.Patronymic,
			Age:         req.Age,
			Gender:      req.Gender,
			Nationality: req.Nationality,
		}

		err = personUpdater.UpdatePerson(id, person)
		if err != nil {
			log.Info("error while updating person", slog.Int("id", id))
			render.JSON(w, r, resp.Error("error while updating person"))
			return
		}

		log.Info("person updated", slog.Int("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
