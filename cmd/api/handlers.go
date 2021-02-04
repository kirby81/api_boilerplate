package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/kirby81/api-boilerplate/internal/todo"
)

func (a *api) addTodo() http.HandlerFunc {
	var req todo.AddTodoRequest
	var resp todo.AddTodoResponse
	return func(w http.ResponseWriter, r *http.Request) {
		req = todo.AddTodoRequest{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			a.infoLog.Printf("Failed to decode the request body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var err error
		resp, err = a.todoService.AddTodo(req)
		if err != nil {
			a.errorLog.Printf("[SERVICE FAILED] Failed to add todo: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			a.errorLog.Printf("Failed to encode addTodo response: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) getTodo() http.HandlerFunc {
	var resp todo.TodoResponse
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			a.infoLog.Printf("Failed to get the id query param: %v", err)
		}

		resp, err = a.todoService.GetTodo(id)
		if err != nil {
			a.errorLog.Printf("[SERVICE FAILED] Failed to get todo: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			a.errorLog.Printf("Failed to encode getTodo response: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) getTodos() http.HandlerFunc {
	var resp todo.GetTodosResponse
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		resp, err = a.todoService.GetTodos()
		if err != nil {
			a.errorLog.Printf("[SERVICE FAILED] Failed to get todo: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			a.errorLog.Printf("Failed to encode getTodo response: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
}

func (a *api) updateTodo() http.HandlerFunc {
	var req todo.UpdateTodoRequest
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			a.infoLog.Printf("Failed to get the id query param: %v", err)
		}

		req = todo.UpdateTodoRequest{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			a.infoLog.Printf("Failed to decode the request body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = a.todoService.UpdateTodo(id, req)
		if err != nil {
			a.errorLog.Printf("[SERVICE FAILED] Failed to add todo: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *api) deleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			a.infoLog.Printf("Failed to get the id query param: %v", err)
		}

		err = a.todoService.DeleteTodo(id)
		if err != nil {
			a.errorLog.Printf("[SERVICE FAILED] Failed to delete todo: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
