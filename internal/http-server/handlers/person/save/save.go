package save

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"people-service/internal/domain/models"
	resp "people-service/internal/lib/api/response"
	"people-service/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Response struct {
	resp.Response
	Id int `json:"id,omitempty"`
}

type PersonSaver interface {
	SavePerson(people models.Person) (id int, err error)
}

type AgeGetter interface {
	GetAge(name string) (age int, err error)
}

type GenderGetter interface {
	GetGender(name string) (gender string, err error)
}

type NationalityGetter interface {
	GetNationality(name string) (nationality string, err error)
}

func New(log *slog.Logger,
	personSaver PersonSaver,
	ageGetter AgeGetter,
	genderGetter GenderGetter,
	nationalityGetter NationalityGetter,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Debug("entered to save person handler")

		var req Request

		err := render.DecodeJSON(r.Body, &req)

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

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", err)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		log.Debug("req is OK")

		person := models.Person{Name: req.Name, Surname: req.Surname}
		if req.Patronymic != "" {
			person.Patronymic = req.Patronymic
		}

		age, err := ageGetter.GetAge(person.Name)
		if err != nil {
			log.Error("failed to get age", err)
		} else {
			person.Age = age
		}

		log.Debug(fmt.Sprintf("age is: %s", age))

		gender, err := genderGetter.GetGender(person.Name)
		if err != nil {
			log.Error("failed to get gender", err)
		} else {
			person.Gender = gender
		}
		log.Debug(fmt.Sprintf("gender is: %s", gender))

		nationality, err := nationalityGetter.GetNationality(person.Name)
		if err != nil {
			log.Error("failed to get nationality", err)
		} else {
			person.Nationality = nationality
		}
		log.Debug(fmt.Sprintf("nationality is: %s", nationality))

		id, err := personSaver.SavePerson(person)

		if errors.Is(err, storage.ErrPersonExists) {
			log.Info("person already exists", slog.String("name", person.Name), slog.String("surname", person.Surname))
			render.JSON(w, r, resp.Error("person already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add person", err)
			render.JSON(w, r, resp.Error("failed to add person"))
			return
		}

		log.Info("person added", slog.Int("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
