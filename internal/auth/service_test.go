package auth

import (
	"errors"
	"testing"

	repomock "github.com/kirby81/api-boilerplate/internal/auth/mock"
	"github.com/kirby81/api-boilerplate/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	mockRepo := repomock.MockRepository{}

	// Default
	service, err := NewService(&mockRepo, "tokensecret")
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// With no repo
	_, err = NewService(nil, "tokensecret")
	assert.Error(t, err)

	// With no secret
	_, err = NewService(&mockRepo, "")
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	mockRepo := repomock.MockRepository{}

	service, err := NewService(&mockRepo, "tokensecret")
	assert.Nil(t, err)

	// Default
	mockRepo.On("GetUser", "foobar@bar.com").Return(&entities.User{
		Email:    "foobar@bar.com",
		Password: "$2y$12$lPEpi.XBHCrzfs.Z0ewCWe62caFEodFFfmtscBqICO47i0cAmY2gS",
	}, nil)
	token, err := service.Login("foobar@bar.com", "password")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// GetUser fail
	mockRepo.On("GetUser", "unknown@user.com").Return(nil, errors.New("error"))
	token, err = service.Login("unknown@user.com", "password")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestSignin(t *testing.T) {
	mockRepo := repomock.MockRepository{}

	service, err := NewService(&mockRepo, "tokensecret")
	assert.Nil(t, err)

	// Default
	mockRepo.On("GetUser", "foobar@bar.com").Return(nil, UserNotFound)
	mockRepo.On("AddUser", "foobar@bar.com", mock.Anything)
	err = service.Signin("foobar@bar.com", "password")
	assert.NoError(t, err)

	// User already exist
	mockRepo.On("GetUser", "existing@user.com").Return(&entities.User{}, nil)
	err = service.Signin("existing@user.com", "password")
	assert.EqualError(t, err, "user already exist")

	// GetUser fail
	mockRepo.On("GetUser", "unknown@user.com").Return(nil, errors.New("error"))
	err = service.Signin("unknown@user.com", "password")
	assert.Error(t, err)
}

func TestParse(t *testing.T) {
	mockRepo := repomock.MockRepository{}

	service, err := NewService(&mockRepo, "tokensecret")
	assert.Nil(t, err)

	// Default
	err = service.Parse(
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1ODUyNDMyMDcsImV4cCI6NDEwOTc2NDk2NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.nnV9ngIQXV5QSJJ-lD9gjyv4DFVuTHOA-5IvPYFo1cQ",
	)
	assert.Nil(t, err)

	// Invalid JWT
	err = service.Parse(
		"foobar",
	)
	assert.Error(t, err)
}
