package query

import (
	"context"
)

type CategoryWithNameDTO struct {
	Id   string
	Name string
}

type CategoryQueryService interface {
	FindAllCategories(ctx context.Context) ([]CategoryWithNameDTO, error)
	FindCategoriesIdIncludes(ctx context.Context, ids []string) ([]CategoryWithNameDTO, error)
}
