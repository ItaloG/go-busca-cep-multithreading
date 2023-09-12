package infra

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ItaloG/go-busca-cep-multithreading/util"
)

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func GetCepFromBrasilApi(cep string, ch chan *util.Response) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}

	if resp.StatusCode == 429 {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}
	var b BrasilApi
	err = json.Unmarshal(body, &b)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}

	ch <- &util.Response{
		From:     "BrasilApi",
		Cep:      b.Cep,
		State:    b.State,
		City:     b.City,
		District: b.Neighborhood,
		Error:    false,
	}
	return
}
