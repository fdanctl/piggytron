// Package appexpensecategory implements the expense category use cases
// (create, find by id, list by user) for the needs / wants / savings taxonomy.
package appexpensecategory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the expense category use cases.
type Service struct {
	repo expensecategory.Repository
	db   *sql.DB
}

// NewService wires the expense category service to its repository.
func NewService(repo expensecategory.Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

// CreateCategory persists a new expense category, rejecting duplicates.
func (s *Service) CreateCategory(
	ctx context.Context,
	userID string,
	name string,
	expenseType string,
) (*expensecategory.ExpenseCategory, error) {
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s is not a valid expense type", et),
			fmt.Errorf("%s is not a valid expense type: %w", et, err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	id, err := util.NewID[expensecategory.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	category, err := expensecategory.New(id, uid, name, et)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, expensecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"An expense category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appexpensecategory.CreateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appexpensecategory.CreateCategory",
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
	expenseType string,
) error {
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appexpensecategory.UpdateCategory",
		)
		return err
	}

	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appexpensecategory.UpdateCategory",
		)
		return err
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s is not a valid expense type", et),
			fmt.Errorf("%s is not a valid expense type: %w", et, err),
			"appexpensecategory.UpdateCategory",
		)
		return err
	}

	category, err := expensecategory.New(cid, uid, name, et)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appexpensecategory.UpdateCategory",
		)
		return err
	}

	err = s.repo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, expensecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"An expense category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appexpensecategory.UpdateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appexpensecategory.UpdateCategory",
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
	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appexpensecategory.Delete",
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

	btx := postgres.NewBudgetRepository(tx)
	budgets, err := btx.FindAllByCategory(ctx, budget.ID(id))
	for _, b := range budgets {
		if b.Amount() > 0 {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf(
					"Can't delete a category with historical data. Budget on %s",
					b.Month().Label(),
				),
				errors.New("can't delete a category with historical data"),
				"appexpensecategory.Delete",
			)
			return err
		}
	}

	etx := postgres.NewExpenseCategoryRepository(tx)
	err = etx.Delete(ctx, cid)
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
) (*expensecategory.ExpenseCategory, error) {
	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}

	cat, err := s.repo.FindByID(ctx, cid)
	if err != nil {
		if errors.Is(err, expensecategory.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The category does not exists",
				fmt.Errorf("failed to found category '%s': %w", id, err),
				"appexpensecategory.FindCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", id, err),
				"appexpensecategory.FindCategory",
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
				expensecategory.ErrNotFound,
			),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}

	return cat, nil
}

// FindAllUserCategories lists the user's expense categories.
func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*expensecategory.ExpenseCategory, error) {
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appexpensecategory.FindAllUserCategories",
		)
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' expense categories: %w", uid, err),
			"appexpensecategory.FindAllUserCategories",
		)
		return nil, err
	}

	return categories, nil
}

func (s *Service) Archive(
	ctx context.Context,
	id string,
) error {
	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appexpensecategory.Archive",
		)
		return err
	}
	return s.repo.Archive(ctx, cid)
}
