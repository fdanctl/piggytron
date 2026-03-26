package query

import (
	"context"
)

type CategoryWithNameDTO struct {
	ID   string
	Name string
}

type CategoryQueryService interface {
	FindAllCategories(ctx context.Context, uid string) ([]CategoryWithNameDTO, error)
	FindCategoriesIDIncludes(ctx context.Context, ids []string) ([]CategoryWithNameDTO, error)
}
