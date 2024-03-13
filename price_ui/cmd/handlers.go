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

func (app *application) handleStorageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var storage struct {
		Baseline  string            `json:"baseline"`
		Discounts map[string]string `json:"discounts"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	} else if err = json.Unmarshal(body, &storage); err != nil {
		app.serverError(w, err)
		return
	}

	params := url.Values{}
	params.Add("data_base_name", storage.Baseline)

	reqURL := app.server_addr + "/change_storage?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса: ", err)
		return
	}
	defer resp.Body.Close()

	for _, matr := range storage.Discounts {
		// Подготовка параметров запроса
		params = url.Values{}
		params.Add("data_base_id", matr)

		// Добавление параметров к URL
		reqURL := app.server_addr + "/get_matrix?" + params.Encode()
		resp, err := http.Get(reqURL)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса: ", err)
			return
		}
		defer resp.Body.Close()
	}

	app.infoLog.Printf("Save storage request: %v, %v", storage.Baseline, storage.Discounts)
	w.WriteHeader(http.StatusOK)
}

func (app *application) handleUpdateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var resp struct {
		Matrix  string  `json:"update_matrix"`
		Updates [][]int `json:"updates"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	} else if err = json.Unmarshal(body, &resp); err != nil {
		app.serverError(w, err)
		return
	}

	_, err = strconv.Atoi(resp.Matrix)
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, tuple := range resp.Updates {
		// Подготовка параметров запроса
		params := url.Values{}
		params.Add("location_id", strconv.Itoa(tuple[0]))
		params.Add("microcategory_id", strconv.Itoa(tuple[1]))
		params.Add("data_base_id", resp.Matrix)
		params.Add("price", strconv.Itoa(tuple[2]))

		// Добавление параметров к URL
		reqURL := app.server_addr + "/set_price?" + params.Encode()
		resp, err := http.Get(reqURL)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса: ", err)
			return
		}
		defer resp.Body.Close()
	}

	app.infoLog.Printf("Save storage request: %v, %v", 1, resp.Updates)
	w.WriteHeader(http.StatusOK)
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
	app.infoLog.Println("addRequestHandler ", data)

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
	var responseObj struct {
		IdMatrix int `json:"id_matrix"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseObj); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	if responseObj.IdMatrix == 0 {
		responseData := struct {
			Result string `json:"res"`
		}{"matrix not found"}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	} else {
		// Отправляем ответ в формате JSON
		responseData := struct {
			Result int `json:"res"`
		}{responseObj.IdMatrix}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

func (app *application) getRequestHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	app.infoLog.Println("getRequestHandler ", data)

	// Подготовка параметров запроса
	params := url.Values{}
	params.Add("data_base_id", data)

	// Добавление параметров к URL
	reqURL := app.server_addr + "/get_matrix?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса: ", err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	var responseObj struct {
		NameMatrix string `json:"matrix_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseObj); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	// Отправляем ответ в формате JSON
	responseData := struct {
		Result string `json:"res"`
	}{responseObj.NameMatrix}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
