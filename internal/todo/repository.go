package todo

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Insert(todo *Todo) error
	GetByID(id int) (*Todo, error)
	GetAll() ([]*Todo, error)
	Update(todo *Todo) error
	Delete(id int) error
}

type PostgreRepository struct {
	DB *sql.DB
}

func NewPostgreRepository(db *sql.DB) Repository {
	return &PostgreRepository{
		DB: db,
	}
}

func (pr *PostgreRepository) Insert(todo *Todo) error {
	insertStmt := `INSERT INTO todos (title, description, done) VALUES ($1, $2, $3) RETURNING id`

	err := pr.DB.QueryRow(insertStmt, todo.title, todo.desc, todo.done).Scan(&todo.id)
	if err != nil {
		return fmt.Errorf("failed to insert todo: %w", err)
	}

	return nil
}

func (pr *PostgreRepository) GetByID(id int) (*Todo, error) {
	stmt := `SELECT id, title, description, done, created_at, updated_at, completed_at FROM todos WHERE id = $1`
	todo := &Todo{}

	err := pr.DB.QueryRow(stmt, id).Scan(
		&todo.id,
		&todo.title,
		&todo.desc,
		&todo.done,
		&todo.createdAt,
		&todo.updatedAt,
		&todo.completedAt,
	)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (pr *PostgreRepository) GetAll() ([]*Todo, error) {
	stmt := `SELECT id, title, description, done, created_at, updated_at, completed_at FROM todos`
	todos := []*Todo{}

	rows, err := pr.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var todo *Todo
	for rows.Next() {
		todo = new(Todo)

		err = rows.Scan(
			&todo.id,
			&todo.title,
			&todo.desc,
			&todo.done,
			&todo.createdAt,
			&todo.updatedAt,
			&todo.completedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (pr *PostgreRepository) Update(todo *Todo) error {
	stmt := `UPDATE todos SET title = $1, description = $2, done = $3, updated_at = $4 WHERE id = $5`

	_, err := pr.DB.Exec(stmt, todo.title, todo.desc, todo.done, todo.updatedAt, todo.id)
	if err != nil {
		return fmt.Errorf("failed to insert todo: %w", err)
	}

	return nil
}

func (pr *PostgreRepository) Delete(id int) error {
	stmt := `DELETE FROM todos WHERE id = $1`

	_, err := pr.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
