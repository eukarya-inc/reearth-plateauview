package plateauapi

import (
	"context"
)

// area

func (*areaResolver[T]) Datasets(ctx context.Context, obj T, input DatasetForAreaQuery) ([]Dataset, error) {
	panic("implement me")
}

// dataset

func (*datasetResolver[T]) Area(ctx context.Context, obj T) (Area, error) {
	panic("implement me")
}

// queries

func (r *queryResolver) Prefecture(ctx context.Context, code string) (*Prefecture, error) {
	panic("implement me")
}

func (r *queryResolver) Prefectures(ctx context.Context, datasetType *DatasetType) ([]*Prefecture, error) {
	panic("implement me")
}

func (r *queryResolver) Municipality(ctx context.Context, code string) (*Municipality, error) {
	panic("implement me")
}

func (r *queryResolver) Municipalities(ctx context.Context, datasetType *DatasetType, prefectureCode *string) ([]*Municipality, error) {
	panic("implement me")
}

func (r *queryResolver) Dataset(ctx context.Context, datasetID string) (Dataset, error) {
	panic("implement me")
}

func (r *queryResolver) Datasets(ctx context.Context, input DatasetQuery) ([]Dataset, error) {
	panic("implement me")
}
