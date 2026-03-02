package user

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, name, password string) error {
	// TODO hash password
	u, err := user.New(user.ID(uuid.New().String()), name, password)
	if err != nil {
		return err
	}

	return s.repo.Save(ctx, u)
}
