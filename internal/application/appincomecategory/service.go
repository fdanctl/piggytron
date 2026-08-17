// Package appincomecategory implements the income category use cases (create,
// find by id, list by user).
package appincomecategory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the income category use cases.
type Service struct {
	repo incomecategory.Repository
	db   *sql.DB
}

// NewService wires the income category service to its repository.
func NewService(repo incomecategory.Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

// CreateCategory persists a new income category, rejecting duplicates.
func (s *Service) CreateCategory(
	ctx context.Context,
	userID string,
	name string,
) (*incomecategory.IncomeCategory, error) {
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	id, err := util.NewID[incomecategory.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	category, err := incomecategory.New(id, uid, name)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, incomecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"An income category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appincomecategory.CreateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appincomecategory.CreateCategory",
			)
		}
		return nil, err
	}
	return category, nil
}

func (s *Service) UpdateCategory(
	ctx context.Context,
	id string,
	userID string,
	name string,
) error {
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appincomecategory.UpdateCategory",
		)
		return err
	}

	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appincomecategory.UpdateCategory",
		)
		return err
	}

	category, err := incomecategory.New(cid, uid, name)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appincomecategory.UpdateCategory",
		)
		return err
	}

	err = s.repo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, incomecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"An income category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appincomecategory.UpdateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appincomecategory.UpdateCategory",
			)
		}
		return err
	}
	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	uid string,
) error {
	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appincomecategory.Delete",
		)
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appincomecategory.Delete",
		)
		return err
	}
	defer tx.Rollback()

	qtx := postgres.NewLedgerQueryService(tx)
	filters := query.NewLedgerFilters(nil, nil, []string{id}, "", "", "", "")
	count, err := qtx.CountFilteredResults(ctx, uid, filters)
	if count > 0 {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Can't delete a category with historical data",
			errors.New("can't delete a category with historical data"),
			"appincomecategory.Delete",
		)
		return err
	}

	itx := postgres.NewIncomeCategoryRepository(tx)
	err = itx.Delete(ctx, cid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to delete '%s' account: %w", id, err),
			"appincomecategory.Delete",
		)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appincomecategory.Delete",
		)
		return err
	}

	return nil
}

// FindCategory returns a category owned by userID, mapping not-found and
// ownership mismatches to KindNotFound.
func (s *Service) FindCategory(
	ctx context.Context,
	id string,
	userID string,
) (*incomecategory.IncomeCategory, error) {
	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}

	cat, err := s.repo.FindByID(ctx, cid)
	if err != nil {
		if errors.Is(err, incomecategory.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The category does not exists",
				fmt.Errorf("failed to found category '%s': %w", id, err),
				"appincomecategory.FindCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", id, err),
				"appincomecategory.FindCategory",
			)
		}
		return nil, err
	}

	if cat.UserID() != uid {
		err = errs.NewAppError(
			errs.KindNotFound,
			"The category does not exists",
			fmt.Errorf(
				"the category does not belong to user '%s': %w",
				uid,
				incomecategory.ErrNotFound,
			),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}

	return cat, nil
}

// FindAllUserCategories lists the user's income categories.
func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*incomecategory.IncomeCategory, error) {
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appincomecategory.FindAllUserCategories",
		)
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' income categories: %w", uid, err),
			"appincomecategory.FindAllUserCategories",
		)
		return nil, err
	}

	return categories, nil
}

func (s *Service) Archive(
	ctx context.Context,
	id string,
) error {
	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appincomecategory.Archive",
		)
		return err
	}
	return s.repo.Archive(ctx, cid)
}
