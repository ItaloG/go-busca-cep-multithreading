package infra

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ItaloG/go-busca-cep-multithreading/util"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetCepFromViaCep(cep string, ch chan *util.Response) {
	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		ch <- &util.Response{From: "", Cep: "", State: "", City: "", District: "", Error: true}
		return
	}

	ch <- &util.Response{
		From:     "ViaCep",
		Cep:      c.Cep,
		State:    c.Logradouro,
		City:     c.Localidade,
		District: c.Bairro,
		Error:    false,
	}
	return
}
