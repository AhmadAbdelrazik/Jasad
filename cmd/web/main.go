package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"	
	"github.com/AhmadAbdelrazik/jasad/internal/model"
)

type Application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	Jasad *model.JasadModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Server Port address")
	dsn := flag.String("dsn", "ahmad:password@/jasad_db?parseTime=true", "MySql Data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	defer db.Close()

	app := &Application{
		infoLog:  infoLog,
		errorLog: errorLog,
		Jasad: &model.JasadModel{DB: db},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	infoLog.Printf("Connecting to %v ...", *addr)
	srv.ListenAndServe()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}