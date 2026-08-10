// Package appuser implements the user use cases: registration, login and
// logout (including from all devices), changing the name and password, and
// Argon2id password hashing (see PasswordHasher). Sessions are issued through
// auth.SessionManager.
package appuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/auth"
	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the user use cases (register, login, logout, profile
// and password changes) and owns session issuance.
type Service struct {
	repo           user.Repository
	hasher         *PasswordHasher
	sessionManager *auth.SessionManager
}

// NewService wires the user service to its repository, password hasher and
// session manager.
func NewService(
	repo user.Repository, hasher *PasswordHasher, sm *auth.SessionManager,
) *Service {
	return &Service{repo: repo, hasher: hasher, sessionManager: sm}
}

// FindByID returns a user by id.
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

// CreateUser registers a new user and returns a session id for it.
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

	sid, err := s.sessionManager.CreateSession(ctx, string(u.ID()))

	return sid, err
}

// ChangeName updates the display name, rejecting duplicates.
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

// ChangePassword verifies the current password, stores the new hash, revokes
// every other session and returns a fresh session id.
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
			fmt.Errorf("failed to verify password: %w", err),
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

	err = s.sessionManager.RevokeAllSessions(ctx, id)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to log out user: %w", err),
			"appuser.ChangePassword",
		)
		return "", err
	}
	sid, err := s.sessionManager.CreateSession(ctx, string(u.ID()))

	return sid, err
}

// LoginUser verifies name and password and returns a new session id.
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

	sid, err := s.sessionManager.CreateSession(ctx, string(u.ID()))
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating session: %w", err),
			"appuser.LoginUser",
		)
	}

	return sid, err
}

// LogoutUser deletes the current session.
func (s *Service) LogoutUser(ctx context.Context, sessionID string) error {
	err := s.sessionManager.DeleteSession(ctx, sessionID)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to log out user: %w", err),
			"appuser.LogoutUser",
		)
	}
	return err
}

// LogoutUserFromAllDevices revokes every session of the user.
func (s *Service) LogoutUserFromAllDevices(ctx context.Context, userID string) error {
	err := s.sessionManager.RevokeAllSessions(ctx, userID)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to log out user: %w", err),
			"appuser.LogoutUserFromAllDevices",
		)
		return err
	}
	return nil
}
