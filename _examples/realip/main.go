package main

import (
	"fmt"
	extmiddleware "github.com/Ferluci/chi-extra-middleware"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(extmiddleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Your IP is %s", r.RemoteAddr)))
	})

	http.ListenAndServe(":3333", r)
}