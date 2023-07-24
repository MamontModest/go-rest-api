package auth

import (
	"context"
	"errors"
	"github.com/mamontmodest/go-rest-api/internal/entity"
	"log"
)

type Service interface {
	Login(ctx context.Context, login, password string) error
	CreateUser(ctx context.Context, login, password string) error
}

type Identity struct {
	entity.Identity
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Login(ctx context.Context, username, password string) error {
	if err := s.authenticate(ctx, username, password); err != nil {
		return nil
	}
	return errors.New("unauthorized")
}

func (s service) authenticate(ctx context.Context, login, password string) error {
	person := entity.Identity{
		Login:    login,
		Password: password,
	}
	err := s.repo.Login(ctx, person)
	log.Println(err, person)
	if err != nil {
		return err
	}
	return err
}

func (s service) CreateUser(ctx context.Context, login, password string) error {
	person := entity.Identity{
		Login:    login,
		Password: password,
	}
	err := s.repo.CreateUser(ctx, person)
	if err != nil {
		return err
	}
	return nil

}
