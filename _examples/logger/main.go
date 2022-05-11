package main

import (
	"net/http"

	extmiddleware "github.com/ferluci/chi-extra-middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(extmiddleware.ZeroLogger())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3333", r)
}
