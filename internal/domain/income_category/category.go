package incomecategory

import "time"

type ID string

func NewId(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type IncomeCategory struct {
	id        ID
	userId    ID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func New(id ID, userId ID, name string) (*IncomeCategory, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now()

	return &IncomeCategory{
		id:        id,
		userId:    userId,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Rehydrate(id ID, userId ID, name string, createdAt, updatedAt time.Time) *IncomeCategory {
	return &IncomeCategory{
		id:        id,
		userId:    userId,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (ic *IncomeCategory) ID() ID {
	return ic.id
}

func (ic *IncomeCategory) UserId() ID {
	return ic.userId
}

func (ic *IncomeCategory) Name() string {
	return ic.name
}

func (ic *IncomeCategory) CreatedAt() time.Time {
	return ic.createdAt
}

func (ic *IncomeCategory) UpdatedAt() time.Time {
	return ic.updatedAt
}

func (ic *IncomeCategory) ChangeName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	ic.name = name
	ic.updatedAt = time.Now()
	return nil
}
