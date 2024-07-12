package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Server Port address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO/t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR/t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	srv.ListenAndServe()
}
