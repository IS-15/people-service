package nationality

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sort"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	Count   int       `json:"count,omitempty"`
	Name    string    `json:"name,omitempty"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type NationalityService struct {
	log     *slog.Logger
	baseUrl string
}

func New(log *slog.Logger, url string) *NationalityService {
	return &NationalityService{log: log, baseUrl: url}
}

func (a *NationalityService) GetNationality(name string) (string, error) {
	const op = "data-prep.nationality.GetNationality"

	log := a.log.With(
		slog.String("op", op),
	)

	req, err := http.NewRequest("GET", a.baseUrl, nil)
	if err != nil {
		log.Error("cannot form new request")
		return "", err
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error while making request")
		return "", err
	}
	defer resp.Body.Close()

	var nResp Response
	err = json.NewDecoder(resp.Body).Decode(&nResp)
	if err != nil {
		log.Error("cannot decode response")
		return "", err
	}

	if len(nResp.Country) < 1 {
		log.Error("no country id found for the person")
		return "", errors.New("no country id found for the person")
	}

	return getSingleNationality(nResp.Country), nil

}

func getSingleNationality(countries []Country) string {

	sort.Slice(countries[:], func(i, j int) bool {
		return countries[i].Probability > countries[j].Probability
	})

	return countries[0].CountryId
}
