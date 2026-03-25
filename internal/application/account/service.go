package account

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/google/uuid"
)

type Service struct {
	repo account.Repository
}

var (
	ErrDuplicate = errors.New("duplicate bank")
	ErrNotNumber = errors.New("not number")
)

func NewService(repo account.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBank(
	ctx context.Context,
	name string,
	currency string,
) (*account.Account, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, errors.New("nil context")
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, errors.New("not sessionInfo")
	}

	_, err := s.repo.FindBankByNameAndUser(ctx, account.ID(sessionInfo.UserId), name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := account.NewId(uuid.New().String())
	if err != nil {
		return nil, err
	}
	uid, err := account.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	account, err := account.NewBank(id, uid, name, currency)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) CreateGoal(
	ctx context.Context,
	name string,
	currency string,
	targetAmount string,
	targetDate string,
	categoryId string,
) (*account.Account, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, errors.New("nil context")
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, errors.New("not sessionInfo")
	}

	_, err := s.repo.FindGoalByNameAndUser(ctx, account.ID(sessionInfo.UserId), name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := account.NewId(uuid.New().String())
	if err != nil {
		return nil, err
	}
	uid, err := account.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}
	cid, err := account.NewId(categoryId)
	if err != nil {
		return nil, err
	}

	tAmount, err := strconv.Atoi(targetAmount)
	if err != nil {
		return nil, ErrNotNumber
	}

	tDate, err := time.Parse(time.DateOnly, targetDate)
	var pDate *time.Time
	if err == nil {
		pDate = &tDate
	}
	account, err := account.NewGoal(id, uid, name, currency, tAmount, pDate, cid)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) ReadOneById(ctx context.Context, id string) (*account.Account, error) {
	newId, err := account.NewId(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindById(ctx, newId)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
) ([]*account.Account, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := account.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) ReadAllBanksByUser(
	ctx context.Context,
) ([]*account.Account, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := account.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllBanksByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) ReadAllGoalsByUser(
	ctx context.Context,
) ([]*account.Account, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := account.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllGoalsByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) ReadIdNamesIncludes(
	ctx context.Context,
	ids []string,
) ([]*account.AccountIdName, error) {
	if len(ids) <= 0 {
		return nil, nil
	}
	return s.repo.FindIdNamesIncludes(ctx, ids)
}
