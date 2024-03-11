package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.handleSimpleRequest("index.tmpl", "/"))
	mux.HandleFunc("/save", app.handleStorageRequest)
	mux.HandleFunc("/metrics", app.handleMetricsRequest)

	fileServer := http.FileServer(http.Dir("../src/"))
	mux.Handle("/src", http.NotFoundHandler())
	mux.Handle("/src/", http.StripPrefix("/src", fileServer))

	return mux
}
