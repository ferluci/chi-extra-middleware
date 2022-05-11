package main

import (
	"net/http"

	extmiddleware "github.com/ferluci/chi-extra-middleware"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := chi.NewRouter()
	// You can use MetricsWithConfig() method for custom configuration
	r.Use(extmiddleware.Metrics())

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("try: curl http://127.0.0.1:3333/metrics"))
	})

	http.ListenAndServe(":3333", r)
}
