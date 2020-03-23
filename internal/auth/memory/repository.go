package memory

import (
	"fmt"

	"github.com/kirby81/api-boilerplate/internal/auth"
	"github.com/kirby81/api-boilerplate/internal/entities"
)

type Repository struct {
	users []entities.User
}

func NewRepository() auth.Repository {
	return &Repository{
		users: make([]entities.User, 0),
	}
}

func (r *Repository) AddUser(email, password string) {
	user := entities.User{
		Email:    email,
		Password: password,
	}

	r.users = append(r.users, user)
	fmt.Println(r.users)
}

func (r *Repository) GetUser(email string) (*entities.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, auth.UserNotFound
}
