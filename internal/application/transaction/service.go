package transaction

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
)

type Service struct {
	repo transaction.Repository
}

func NewService(r transaction.Repository) *Service {
	return &Service{repo: r}
}

// create income
// create expense
// create transfer

func (s *Service) ReadOneByID(ctx context.Context, id string) (*transaction.Transaction, error) {
	newID, err := transaction.NewID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, newID)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
	userID string,
	page uint,
) ([]*transaction.Transaction, error) {
	uid, err := transaction.NewID(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) ReadAllByAccount(
	ctx context.Context,
	aid string,
) ([]*transaction.Transaction, error) {
	newID, err := transaction.NewID(aid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByAccount(ctx, newID)
}

func (s *Service) ReadAllByCategory(
	ctx context.Context,
	cid string,
) ([]*transaction.Transaction, error) {
	newID, err := transaction.NewID(cid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByCategory(ctx, newID)
}

// func (s *Service) ReadFiltered(
// 	ctx context.Context,
// 	filters *query.TransactionFilters,
// 	page uint,
// ) ([]query.TransactionDTO, bool, error) {
// 	v := ctx.Value(middleware.UserKey)
// 	if v == nil {
// 		return nil, false, nil
// 	}
//
// 	sessionInfo, ok := v.(*rdb.SessionInfo)
// 	if !ok {
// 		return nil, false, nil
// 	}
//
// 	transactions, err := s.query.FindFiltered(
// 		ctx,
// 		sessionInfo.UserID,
// 		filters,
// 		LIMIT+1,
// 		LIMIT*page-LIMIT,
// 	)
// 	if err != nil {
// 		return nil, false, err
// 	}
//
// 	var hasMore bool
// 	if len(transactions) == LIMIT+1 {
// 		hasMore = true
// 		transactions = transactions[0 : len(transactions)-1]
// 	}
//
// 	return transactions, hasMore, nil
// }
//
// func (s *Service) ReadFilteredWithCount(
// 	ctx context.Context,
// 	filters *query.TransactionFilters,
// 	page uint,
// ) (*query.TransactionsWithTotalCount, bool, error) {
// 	v := ctx.Value(middleware.UserKey)
// 	if v == nil {
// 		return nil, false, nil
// 	}
//
// 	sessionInfo, ok := v.(*rdb.SessionInfo)
// 	if !ok {
// 		return nil, false, nil
// 	}
//
// 	tWithCount, err := s.query.FindFilteredWithCount(
// 		ctx,
// 		sessionInfo.UserID,
// 		filters,
// 		LIMIT+1,
// 		LIMIT*page-LIMIT,
// 	)
// 	if err != nil {
// 		return nil, false, nil
// 	}
//
// 	var hasMore bool
// 	if len(tWithCount.Data) == LIMIT+1 {
// 		hasMore = true
// 		tWithCount.Data = tWithCount.Data[0 : len(tWithCount.Data)-1]
// 	}
//
// 	return tWithCount, hasMore, nil
// }
//
// func (s *Service) CountFilteredResults(
// 	ctx context.Context,
// 	filters *query.TransactionFilters,
// ) (int, error) {
// 	v := ctx.Value(middleware.UserKey)
// 	if v == nil {
// 		return 0, nil
// 	}
//
// 	sessionInfo, ok := v.(*rdb.SessionInfo)
// 	if !ok {
// 		return 0, nil
// 	}
//
// 	return s.query.CountFilteredResults(ctx, sessionInfo.UserID, filters)
// }
