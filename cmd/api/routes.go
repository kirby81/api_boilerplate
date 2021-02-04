package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (a *api) routes() http.Handler {
	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./web/static/"))

	r.Use(a.logRequest)

	r.Get("/", a.home())
	r.Route("/api/todos", func(r chi.Router) {
		r.Get("/", a.getTodos())
		r.Post("/", a.addTodo())

		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", a.getTodo())
			r.Put("/", a.updateTodo())
			r.Delete("/", a.deleteTodo())
		})
	})
	r.Handle("/static", http.StripPrefix("/static", fileServer))

	return r
}
