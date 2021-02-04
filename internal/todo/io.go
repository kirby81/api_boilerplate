package todo

import (
	"time"
)

type AddTodoRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Done  bool   `json:"done"`
}

type AddTodoResponse struct {
	ID int `json:"id"`
}

type UpdateTodoRequest struct {
	Title *string `json:"title,omitempty"`
	Desc  *string `json:"desc,omitempty"`
	Done  *bool   `json:"done,omitempty"`
}

type TodoResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetTodosResponse struct {
	Todos []TodoResponse `json:"todos"`
}
