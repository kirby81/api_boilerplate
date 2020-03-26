package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	tokenSecret string
	repo        Repository
}

func NewService(repo Repository, secret string) (*Service, error) {
	if secret == "" {
		return nil, errors.New("token secret can't be empty")
	}

	if repo == nil {
		return nil, errors.New("no repository provided")
	}

	return &Service{
		tokenSecret: secret,
		repo:        repo,
	}, nil
}

func (s *Service) Login(email string, password string) (string, error) {
	user, err := s.repo.GetUser(email)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	// Hash compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("failed to compare password: %w", err)
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	signedToken, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) Signin(email string, password string) error {
	_, err := s.repo.GetUser(email)
	if err != nil {
		if !errors.Is(err, UserNotFound) {
			return fmt.Errorf("failed to get user: %w", err)
		}
	} else {
		return errors.New("user already exist")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	s.repo.AddUser(email, string(hash))

	return nil
}

func (s *Service) Parse(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.tokenSecret), nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	return nil
}
