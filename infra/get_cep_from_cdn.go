package infra

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ItaloG/go-busca-cep-multithreading/util"
)

type CdnCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func GetCepFromCdn(cep string, ch chan *util.Response) {
	resp, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
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
	var c CdnCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}

	ch <- &util.Response{
		From:     "CDN",
		Cep:      c.Code,
		State:    c.State,
		City:     c.Address,
		District: c.District,
		Error:    false,
	}
	return
}
