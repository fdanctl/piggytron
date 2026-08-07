package appuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/fdanctl/piggytron/internal/errs"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/util"
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

func (s *Service) FindByID(ctx context.Context, id string) (*user.User, error) {
	uid, err := util.ParseID[user.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appuser.FindByID",
		)
		return nil, err
	}

	u, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"User not found",
				fmt.Errorf("failed find user '%s': %w", id, err),
				"appuser.FindByID",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed finding user: %w", err),
				"appuser.FindByID",
			)
		}
		return nil, err
	}
	return u, nil
}

func (s *Service) CreateUser(ctx context.Context, name, password string) (string, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed hashing password: %w", err),
			"appuser.CreateUser",
		)
		return "", err
	}

	id, err := util.NewID[user.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appuser.CreateUser",
		)
		return "", err
	}

	u, err := user.New(id, name, hash)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"",
			fmt.Errorf("failed creating user: %w", err),
			"appuser.CreateUser",
		)
		return "", err
	}

	err = s.repo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, user.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindConflict,
				"User already exists",
				fmt.Errorf("failed saving user '%s': %w", u.Name(), err),
				"appuser.CreateUser",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving user: %w", err),
				"appuser.CreateUser",
			)
		}
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
		UserID: string(u.ID()), SessionVersion: 1,
	})

	return sid, err
}

func (s *Service) ChangeName(ctx context.Context, id, name string) error {
	uid, err := util.ParseID[user.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appuser.ChangeName",
		)
		return err
	}

	if err := user.ValidateName(name); err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Invalid name",
			fmt.Errorf("failed to validate name '%s': %w", name, err),
			"appuser.ChangeName",
		)
		return err
	}

	err = s.repo.UpdateName(ctx, uid, name)
	if err != nil {
		if errors.Is(err, user.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindConflict,
				"User already exists",
				fmt.Errorf("failed updating user '%s': %w", name, err),
				"appuser.ChangeName",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed updatting user: %w", err),
				"appuser.ChangeName",
			)
		}
		return err
	}

	return nil
}

func (s *Service) ChangePassword(
	ctx context.Context,
	id, currPassword, newPassword string,
) (string, error) {
	u, err := s.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	match, err := s.hasher.Verify(u.PasswordHash(), currPassword)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed verifing password: %w", err),
			"appuser.ChangePassword",
		)
		return "", err
	}
	if !match {
		err = errs.NewAppError(
			errs.KindValidation,
			"Incorrect password",
			fmt.Errorf("failed to verify password: %w", user.ErrWrongPassword),
			"appuser.ChangePassword",
		)
		return "", err
	}

	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed hashing password: %w", err),
			"appuser.ChangePassword",
		)
		return "", err
	}

	if err := u.ChangePassword(hash); err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			"Invalid new password",
			fmt.Errorf("invalid new password: %w", err),
			"appuser.ChangePassword",
		)
		return "", err
	}

	if err := s.repo.Update(ctx, u); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving user: %w", err),
			"appuser.ChangePassword",
		)
		return "", err
	}

	sid, err := s.sessionStore.Set(ctx, &rdb.SessionInfo{
		UserID: string(u.ID()), SessionVersion: 1,
	})

	return sid, err
}

func (s *Service) LoginUser(ctx context.Context, name, password string) (string, error) {
	u, err := s.repo.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindValidation,
				"Name or password are invalid",
				fmt.Errorf("failed finding user '%s': %w", name, err),
				"appuser.LoginUser",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed finding user '%s': %w", name, err),
				"appuser.LoginUser",
			)
		}
		return "", err
	}

	match, err := s.hasher.Verify(u.PasswordHash(), password)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed verifing password: %w", err),
			"appuser.LoginUser",
		)
		return "", err
	}
	if !match {
		err = errs.NewAppError(
			errs.KindValidation,
			"Name or password are invalid",
			fmt.Errorf("failed finding user '%s': %w", name, user.ErrWrongPassword),
			"appuser.LoginUser",
		)
		return "", err
	}

	// TODO add session version to pg and pass it instead
	sid, err := s.sessionStore.Set(ctx, &rdb.SessionInfo{
		UserID: string(u.ID()), SessionVersion: 1,
	})
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user: %w", user.ErrWrongPassword),
			"appuser.LoginUser",
		)
	}

	return sid, err
}

func (s *Service) LogoutUser(ctx context.Context, sessionID string) error {
	err := s.sessionStore.Remove(ctx, sessionID)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to log out user: %w", err),
			"appuser.LogoutUser",
		)
	}
	return err
}

func (s *Service) LogoutUserFromAllDevices(ctx context.Context, id, sessionID string) error {
	err := s.sessionStore.Remove(ctx, sessionID)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to log out user: %w", err),
			"appuser.LogoutUserFromAllDevices",
		)
		return err
	}
	// TODO finish
	return nil
}
