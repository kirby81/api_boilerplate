package todo

import (
	"time"
)

type TodoService struct {
	repo Repository
}

func NewTodoService(repo Repository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) AddTodo(req AddTodoRequest) (AddTodoResponse, error) {
	todo := &Todo{
		title: req.Title,
		desc:  req.Desc,
		done:  req.Done,
	}

	err := s.repo.Insert(todo)
	if err != nil {
		return AddTodoResponse{}, err
	}

	resp := AddTodoResponse{
		ID: todo.id,
	}

	return resp, nil
}

func (s *TodoService) GetTodo(id int) (TodoResponse, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return TodoResponse{}, err
	}

	resp := TodoResponse{
		ID:        todo.id,
		Title:     todo.title,
		Desc:      todo.desc,
		Done:      todo.done,
		CreatedAt: todo.createdAt,
	}

	return resp, nil
}

func (s *TodoService) GetTodos() (GetTodosResponse, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return GetTodosResponse{}, err
	}

	resp := GetTodosResponse{
		Todos: []TodoResponse{},
	}
	for _, todo := range todos {
		todoResponse := TodoResponse{
			ID:        todo.id,
			Title:     todo.title,
			Desc:      todo.desc,
			Done:      todo.done,
			CreatedAt: todo.createdAt,
		}
		resp.Todos = append(resp.Todos, todoResponse)
	}

	return resp, nil
}

func (s *TodoService) UpdateTodo(id int, req UpdateTodoRequest) error {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	change := false
	if req.Title != nil && *req.Title != todo.title {
		todo.title = *req.Title
		change = true
	}
	if req.Desc != nil && *req.Desc != todo.desc {
		todo.desc = *req.Desc
		change = true
	}
	if req.Done != nil && *req.Done != todo.done {
		todo.done = *req.Done
		change = true
	}

	if change {
		todo.updatedAt = time.Now()
		err = s.repo.Update(todo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TodoService) DeleteTodo(id int) error {
	err := s.repo.Delete(id)
	return err
}
