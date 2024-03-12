package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.handleSimpleRequest("index.tmpl", "/"))
	mux.HandleFunc("/storage", app.handleStorageRequest)
	mux.HandleFunc("/update", app.handleUpdateRequest)
	mux.HandleFunc("/metrics", app.handleMetricsRequest)
	mux.HandleFunc("/add_request", app.addRequestHandler)
	mux.HandleFunc("/get_request", app.getRequestHandler)

	fileServer := http.FileServer(http.Dir("./src/files"))
	mux.Handle("/files", http.NotFoundHandler())
	mux.Handle("/files/", http.StripPrefix("/files", fileServer))

	return mux
}
