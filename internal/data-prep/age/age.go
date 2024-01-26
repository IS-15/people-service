package age

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	//resp.Response
	Count int    `json:"count,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age" validate:"required"`
}

type AgeService struct {
	log     *slog.Logger
	baseUrl string
}

func New(log *slog.Logger, url string) *AgeService {
	return &AgeService{log: log, baseUrl: url}
}

func (a *AgeService) GetAge(name string) (int, error) {
	const op = "data-prep.age.GetAge"

	log := a.log.With(
		slog.String("op", op),
	)

	req, err := http.NewRequest("GET", a.baseUrl, nil)
	if err != nil {
		log.Error("cannot form new request")
		return 0, err
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()
	//req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error while making request")
		return 0, err
	}
	defer resp.Body.Close()

	var ageResp Response
	err = json.NewDecoder(resp.Body).Decode(&ageResp)
	if err != nil {
		log.Error("cannot decode response")
		return 0, err
	}

	return ageResp.Age, nil

}
