package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (a *api) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(a.logRequest)

	r.Route("/api/todos", func(r chi.Router) {
		r.Get("/", a.getTodos())
		r.Post("/", a.addTodo())

		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", a.getTodo())
			r.Put("/", a.updateTodo())
			r.Delete("/", a.deleteTodo())
		})
	})

	return r
}
