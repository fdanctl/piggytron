// Package user defines the user aggregate: identity, display name and
// Argon2id password hash.
package user

import "time"

// ID is a user identifier.
type ID string

// User is the authenticated identity: display name and Argon2id password hash.
type User struct {
	id           ID
	name         string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
}

// New builds a validated user from a raw (already hashed) password.
func New(id ID, name, passwordHash string) (*User, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}

	now := time.Now()

	return &User{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Rehydrate rebuilds a User from persistence without re-running validation.
func Rehydrate(
	id ID,
	name, passwordHash string,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:           id,
		name:         name,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ValidateName reports whether name is acceptable for a user.
func ValidateName(name string) error {
	if name == "" || len(name) > 50 {
		return ErrInvalidName
	}
	return nil
}

// ID returns the user id.
func (u *User) ID() ID {
	return u.id
}

// Name returns the display name.
func (u *User) Name() string {
	return u.name
}

// PasswordHash returns the Argon2id hash of the user's password.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// CreatedAt returns when the user was created.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ChangeName updates the display name.
func (u *User) ChangeName(name string) error {
	if name == "" || len(name) > 50 {
		return ErrInvalidName
	}
	u.name = name
	u.updatedAt = time.Now()
	return nil
}

// ChangePassword replaces the stored hash with a new Argon2id hash.
func (u *User) ChangePassword(passwordHash string) error {
	if passwordHash == "" {
		return ErrInvalidPassword
	}
	u.passwordHash = passwordHash
	u.updatedAt = time.Now()
	return nil
}
