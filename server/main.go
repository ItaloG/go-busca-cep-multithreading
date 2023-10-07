package main

import (
	"net/http"

	"github.com/ItaloG/go-busca-cep-multithreading/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", handlers.GetCepHandler)

	http.ListenAndServe(":8000", r)
}
