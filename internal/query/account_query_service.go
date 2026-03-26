package query

import "context"

type AccountIDName struct {
	ID   string
	Name string
}

type AccountQueryService interface {
	FindIDNamesIncludes(ctx context.Context, ids []string) ([]AccountIDName, error)
}
