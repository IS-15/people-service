package get

import (
	"net/http"
	"people-service/internal/domain/models"
	resp "people-service/internal/lib/api/response"
	queryparam "people-service/internal/lib/query-param"
	"people-service/internal/lib/routing"

	"log/slog"

	"github.com/go-chi/render"
)

type PersonGetter interface {
	GetPerson(params queryparam.Params) ([]models.Person, error)
}

func New(log *slog.Logger, personGetter PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.person.get.New"

		qParams := queryparam.Params{}

		qParams.Age = r.URL.Query().Get(routing.AgeParam)
		qParams.Gender = r.URL.Query().Get(routing.GenderParam)
		qParams.Id = r.URL.Query().Get(routing.IdParam)
		qParams.Name = r.URL.Query().Get(routing.NameParam)
		qParams.Surname = r.URL.Query().Get(routing.SurnameParam)
		qParams.Patronymic = r.URL.Query().Get(routing.PatronymicParam)
		qParams.Nationality = r.URL.Query().Get(routing.NationalityParam)
		qParams.Offset = r.URL.Query().Get(routing.OffsetParam)
		qParams.Limit = r.URL.Query().Get(routing.LimitParam)

		persons, err := personGetter.GetPerson(qParams)
		if err != nil {
			log.Error("failed to get persons", err)
			render.JSON(w, r, resp.Error("failed to get persons"))
			return
		}

		render.JSON(w, r, persons)

	}
}
