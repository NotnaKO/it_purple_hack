package main

import "log"

type application struct {
	server_addr string
	errorLog    *log.Logger
	infoLog     *log.Logger
}
