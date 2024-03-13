package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
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

	// TODO send request to server
}

type updater struct {
	Matrix  string  `json:"update_matrix"`
	Updates [][]int `json:"updates"`
}

func (app *application) handleUpdateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var resp updater
	body, err := ioutil.ReadAll(r.Body)
	app.errorLog.Println(body)
	if err != nil {
		app.serverError(w, err)
		return
	} else if err = json.Unmarshal(body, &resp); err != nil {
		app.serverError(w, err)
		return
	}

	num, err := strconv.Atoi(resp.Matrix)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Printf("Save storage request: %v, %v", num, resp.Updates)
	w.WriteHeader(http.StatusOK)

	// TODO send request to server
}

func (app *application) handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// TODO send request for metrics
}

func (app *application) addRequestHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	app.infoLog.Println(data)

	// Подготовка параметров запроса
	params := url.Values{}
	params.Add("data_base_name", data)

	// Добавление параметров к URL
	reqURL := app.server_addr + "/get_id?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса: ", err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Отправляем ответ в формате JSON
	responseData := struct {
		Result string `json:"res"`
	}{string(body)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func (app *application) getRequestHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	app.infoLog.Println("Hi ", data)

	// Подготовка параметров запроса
	params := url.Values{}
	params.Add("data_base_name", data)

	// Добавление параметров к URL
	reqURL := app.server_addr + "/get_id?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса: ", err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Отправляем ответ в формате JSON
	responseData := struct {
		Result int `json:"res"`
	}{1}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
