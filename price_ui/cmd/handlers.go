package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	mes := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, mes)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) handleSimpleRequest(fileName string, path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			app.notFound(w)
			return
		}

		ts, err := template.ParseFiles("./src/" + fileName)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
		}
	}
}

type storage struct {
	Baseline  string            `json:"baseline"`
	Discounts map[string]string `json:"discounts"`
}

func (app *application) handleStorageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var resp storage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	} else if err = json.Unmarshal(body, &resp); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("Save storage request: %v, %v", resp.Baseline, resp.Discounts)
	w.WriteHeader(http.StatusOK)

	// TODO save config file
}

func (app *application) handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// TODO send request for metrics

	fmt.Fprint(w, "Gotcha!")
}
