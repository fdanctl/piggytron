package user

import (
	"context"
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/google/uuid"
)

var (
	ErrWrongPassword = errors.New("password not match")
	ErrUserExists    = errors.New("user already taken")
)

type Service struct {
	repo   user.Repository
	hasher *PasswordHasher
}

func NewService(repo user.Repository, hasher *PasswordHasher) *Service {
	return &Service{repo: repo, hasher: hasher}
}

func (s *Service) CreateUser(ctx context.Context, name, password string) error {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	_, err = s.repo.FindByName(ctx, name)
	if err == nil {
		return ErrUserExists
	}
	u, err := user.New(user.ID(uuid.New().String()), name, hash)
	if err != nil {
		return err
	}

	return s.repo.Save(ctx, u)
}

func (s *Service) LoginUser(ctx context.Context, name, password string) error {
	u, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return err
	}

	match, err := s.hasher.Verify(u.PasswordHash(), password)
	if err != nil {
		return err
	}
	if !match {
		return ErrWrongPassword
	}

	return nil
}
