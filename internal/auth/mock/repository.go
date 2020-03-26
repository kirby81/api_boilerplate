package mock

import (
	"github.com/kirby81/api-boilerplate/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mr *MockRepository) AddUser(email, password string) {
	mr.Called(email, password)
}

func (mr *MockRepository) GetUser(email string) (*entities.User, error) {
	args := mr.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}
