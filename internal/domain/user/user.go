package user

import "time"

type ID string

func NewID(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type User struct {
	id          ID
	name        string
	paswordHash string
	createdAt   time.Time
	updatedAt   time.Time
}

func New(id ID, name, paswordHash string) (*User, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if paswordHash == "" {
		return nil, ErrInvalidPassword
	}

	now := time.Now()

	return &User{
		id:          id,
		name:        name,
		paswordHash: paswordHash,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func Rehydrate(id ID, name, paswordHash string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:          id,
		name:        name,
		paswordHash: paswordHash,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) PasswordHash() string {
	return u.paswordHash
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.name = name
	u.updatedAt = time.Now()
	return nil
}

// TODO change password
