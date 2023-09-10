package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/ItaloG/go-busca-cep-multithreading/infra"
	"github.com/ItaloG/go-busca-cep-multithreading/util"
	"github.com/go-chi/chi/v5"
)

func GetCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Por favor informe um cep")
		return
	}

	isValid, err := regexp.MatchString(`\d{5}-\d{3}`, cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isValid == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Cep com formato inválido. Por favor informe o cep no formato 00000-000")
		return
	}

	cdnCh := make(chan *util.Response)
	viacepCh := make(chan *util.Response)

	go infra.GetCepFromCdn(cep, cdnCh)
	go infra.GetCepFromViaCep(cep, viacepCh)

	for {
		select {
		case msg := <-cdnCh:
			if msg.Error == true {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Falha ao buscar cep. Cep inválido")
				return
			}

			res := fmt.Sprintf(
				"Cep encontrado com sucesso pelo %s. Resultado: cep: %s, estado: %s, cidade: %s, bairro: %s",
				msg.From, msg.Cep, msg.State, msg.City, msg.District,
			)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return

		case msg := <-viacepCh:
			if msg.Error == true {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Falha ao buscar cep. Cep inválido")
				return
			}

			res := fmt.Sprintf(
				"Cep encontrado com sucesso pelo %s. Resultado: cep: %s, estado: %s, cidade: %s, bairro: %s",
				msg.From, msg.Cep, msg.State, msg.City, msg.District,
			)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return

		case <-time.After(time.Second):
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Não foi póssivel encontrar cep!")
			return
		}
	}
}
