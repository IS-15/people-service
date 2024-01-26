package gender

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	Count       int     `json:"count,omitempty"`
	Name        string  `json:"name,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	Probability float32 `json:"probability,omitempty"`
}

type GenderService struct {
	log     *slog.Logger
	baseUrl string
}

func New(log *slog.Logger, url string) *GenderService {
	return &GenderService{log: log, baseUrl: url}
}

func (a *GenderService) GetGender(name string) (string, error) {
	const op = "data-prep.gender.GetGender"

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

	var gResp Response
	err = json.NewDecoder(resp.Body).Decode(&gResp)
	if err != nil {
		log.Error("cannot decode response")
		return "", err
	}

	return gResp.Gender, nil

}
