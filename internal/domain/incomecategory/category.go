// Package incomecategory defines the income category aggregate.
package incomecategory

import "time"

// ID is an income category identifier.
type ID string

// CategoryStatus marks the status of the account (active or archived).
type CategoryStatus string

const (
	ActiveStatus   CategoryStatus = "active"
	ArchivedStatus CategoryStatus = "archived"
)

// IncomeCategory groups income entries, owned by one user.
type IncomeCategory struct {
	id         ID
	userID     ID
	name       string
	status     CategoryStatus
	archivedAt *time.Time
	createdAt  time.Time
	updatedAt  time.Time
}

// New builds a validated income category.
func New(id ID, userID ID, name string) (*IncomeCategory, error) {
	if name == "" || len(name) > 30 {
		return nil, ErrInvalidName
	}

	now := time.Now()

	return &IncomeCategory{
		id:        id,
		userID:    userID,
		name:      name,
		status:    ActiveStatus,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Rehydrate rebuilds an IncomeCategory from persistence without re-running
// validation (the database constraints already guard it).
func Rehydrate(
	id ID,
	userID ID,
	name string,
	status CategoryStatus,
	archivedAt *time.Time,
	createdAt, updatedAt time.Time,
) *IncomeCategory {
	return &IncomeCategory{
		id:         id,
		userID:     userID,
		name:       name,
		status:     status,
		archivedAt: archivedAt,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

// ID returns the category id.
func (ic *IncomeCategory) ID() ID {
	return ic.id
}

// UserID returns the id of the user who owns the category.
func (ic *IncomeCategory) UserID() ID {
	return ic.userID
}

// Name returns the category name.
func (ic *IncomeCategory) Name() string {
	return ic.name
}

// Status returns the category status.
func (ic *IncomeCategory) Status() CategoryStatus {
	return ic.status
}

// ArchivedAt returns the date the category was archived or nil it's active.
func (ic *IncomeCategory) ArchivedAt() *time.Time {
	return ic.archivedAt
}

// CreatedAt returns when the category was created.
func (ic *IncomeCategory) CreatedAt() time.Time {
	return ic.createdAt
}

// UpdatedAt returns when the category was last updated.
func (ic *IncomeCategory) UpdatedAt() time.Time {
	return ic.updatedAt
}

// ChangeName renames the category.
func (ic *IncomeCategory) ChangeName(name string) error {
	if name == "" || len(name) > 30 {
		return ErrInvalidName
	}
	ic.name = name
	ic.updatedAt = time.Now()
	return nil
}
