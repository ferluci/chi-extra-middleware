package main

import (
	"fmt"
	"net/http"

	extmiddleware "github.com/ferluci/chi-extra-middleware"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Use(extmiddleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Your IP is %s", r.RemoteAddr)))
	})

	http.ListenAndServe(":3333", r)
}
