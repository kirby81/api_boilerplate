package auth

import (
	"errors"

	"github.com/kirby81/api-boilerplate/internal/entities"
)

var UserNotFound error = errors.New("user not found")

type Repository interface {
	AddUser(email, password string)
	GetUser(email string) (*entities.User, error)
}
