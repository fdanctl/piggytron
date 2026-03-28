package query

import (
	"context"
)

type CategoryNameDTO struct {
	ID   string
	Name string
}

type CategoryQueryService interface {
	FindAllCategories(ctx context.Context, uid string) ([]CategoryNameDTO, error)
	FindCategoriesIDIncludes(ctx context.Context, ids []string) ([]CategoryNameDTO, error)
}
