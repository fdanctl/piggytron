package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/user"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/google/uuid"
)

var (
	ErrWrongPassword = errors.New("password not match")
	ErrUserExists    = errors.New("user already taken")
)

type Service struct {
	repo         user.Repository
	hasher       *PasswordHasher
	sessionStore *rdb.SessionStore
}

func NewService(
	repo user.Repository, hasher *PasswordHasher, ss *rdb.SessionStore,
) *Service {
	return &Service{repo: repo, hasher: hasher, sessionStore: ss}
}

func (s *Service) CreateUser(ctx context.Context, name, password string) (string, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}

	_, err = s.repo.FindByName(ctx, name)
	if err == nil {
		return "", ErrUserExists
	}
	u, err := user.New(user.ID(uuid.New().String()), name, hash)
	if err != nil {
		return "", err
	}

	err = s.repo.Save(ctx, u)
	if err != nil {
		return "", err
	}

	// TODO add session version to pg and pass it instead
	// session version will for revoke other sessions of the user
	// ex:
	// 1. user updates pwd
	// 2. updates in pg session_vesion + 1
	// 3. create new session with updated version
	// every time a request is made compare the session version with the
	// version on pg if lower session is not valid
	sid, err := s.sessionStore.Set(ctx, &rdb.SessionInfo{
		UserId: string(u.ID()), SessionVersion: 1,
	})

	return sid, err
}

// TODO change name
// TODO change password

func (s *Service) LoginUser(ctx context.Context, name, password string) (string, error) {
	u, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return "", err
	}

	match, err := s.hasher.Verify(u.PasswordHash(), password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", ErrWrongPassword
	}

	// TODO add session version to pg and pass it instead
	sid, err := s.sessionStore.Set(ctx, &rdb.SessionInfo{
		UserId: string(u.ID()), SessionVersion: 1,
	})

	return sid, err
}

func (s *Service) LogoutUser(ctx context.Context) error {
	v := ctx.Value("user")
	if v == nil {
		return nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		fmt.Println("not sessionInfo")
		return nil
	}

	return s.sessionStore.Remove(ctx, sessionInfo.UserId)
}
