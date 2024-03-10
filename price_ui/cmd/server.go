package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

/*
Routes:

/ - ANY - main.page.tmpl
/create - POST - request to send email
/about_us - ANY - about_us.page.tmpl
/delivery - ANY - delivery.page.tmpl
/goods_and_services - ANY - goods_and_services.page.tmpl

*/

func main() {

	addr := flag.String("addr", "127.0.0.1:9992", "Сетевой адрес HTTP")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Запуск веб-сервера на %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
