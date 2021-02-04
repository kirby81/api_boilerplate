package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/kirby81/api-boilerplate/internal/common/http"
	"github.com/kirby81/api-boilerplate/internal/todo"

	_ "github.com/lib/pq"
)

type api struct {
	todoService *todo.TodoService
	infoLog     *log.Logger
	errorLog    *log.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP  network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("postgres", "user=web password=boilerplate dbname=boilerplate sslmode=disable")
	if err != nil {
		errorLog.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	repo := todo.NewPostgreRepository(db)
	api := api{
		todoService: todo.NewTodoService(repo),
		infoLog:     infoLog,
		errorLog:    errorLog,
	}

	srv := http.NewDefaultServer(*addr, api.routes(), errorLog)

	infoLog.Printf("Starting server on %s", *addr)
	if err = srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
