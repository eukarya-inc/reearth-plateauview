package plateauapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

// Area is the resolver for the area field.
func (r *genericDatasetResolver) Area(ctx context.Context, obj *GenericDataset) (Area, error) {
	return to[Area](r.Repo.Node(ctx, obj.AreaID))
}

// Type is the resolver for the type field.
func (r *genericDatasetResolver) Type(ctx context.Context, obj *GenericDataset) (*GenericDatasetType, error) {
	return to[*GenericDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *genericDatasetItemResolver) Parent(ctx context.Context, obj *GenericDatasetItem) (*GenericDataset, error) {
	return to[*GenericDataset](r.Repo.Node(ctx, obj.ParentID))
}

// Datasets is the resolver for the datasets field.
func (r *municipalityResolver) Datasets(ctx context.Context, obj *Municipality, input DatasetForAreaQuery) ([]Dataset, error) {
	return r.Repo.Datasets(ctx, DatasetQuery{
		AreaCodes:    []AreaCode{obj.Code},
		ExcludeTypes: input.ExcludeTypes,
		IncludeTypes: input.IncludeTypes,
		SearchTokens: input.SearchTokens,
	})
}

// Area is the resolver for the area field.
func (r *plateauDatasetResolver) Area(ctx context.Context, obj *PlateauDataset) (Area, error) {
	return to[Area](r.Repo.Node(ctx, obj.AreaID))
}

// Type is the resolver for the type field.
func (r *plateauDatasetResolver) Type(ctx context.Context, obj *PlateauDataset) (*PlateauDatasetType, error) {
	return to[*PlateauDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *plateauDatasetItemResolver) Parent(ctx context.Context, obj *PlateauDatasetItem) (*PlateauDataset, error) {
	return to[*PlateauDataset](r.Repo.Node(ctx, obj.ParentID))
}

// Datasets is the resolver for the datasets field.
func (r *prefectureResolver) Datasets(ctx context.Context, obj *Prefecture, input DatasetForAreaQuery) ([]Dataset, error) {
	return r.Repo.Datasets(ctx, DatasetQuery{
		AreaCodes:    []AreaCode{obj.Code},
		ExcludeTypes: input.ExcludeTypes,
		IncludeTypes: input.IncludeTypes,
		SearchTokens: input.SearchTokens,
	})
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id ID) (Node, error) {
	return r.Repo.Node(ctx, id)
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []ID) ([]Node, error) {
	return r.Repo.Nodes(ctx, ids)
}

// Area is the resolver for the area field.
func (r *queryResolver) Area(ctx context.Context, code AreaCode) (Area, error) {
	return r.Repo.Area(ctx, code)
}

// Areas is the resolver for the areas field.
func (r *queryResolver) Areas(ctx context.Context, input AreaQuery) ([]Area, error) {
	return r.Repo.Areas(ctx, input)
}

// DatasetTypes is the resolver for the datasetTypes field.
func (r *queryResolver) DatasetTypes(ctx context.Context, input DatasetTypeQuery) ([]DatasetType, error) {
	return r.Repo.DatasetTypes(ctx, input)
}

// Datasets is the resolver for the datasets field.
func (r *queryResolver) Datasets(ctx context.Context, input DatasetQuery) ([]Dataset, error) {
	return r.Repo.Datasets(ctx, input)
}

// PlateauSpecs is the resolver for the plateauSpecs field.
func (r *queryResolver) PlateauSpecs(ctx context.Context) ([]*PlateauSpec, error) {
	return r.Repo.PlateauSpecs(ctx)
}

// Years is the resolver for the years field.
func (r *queryResolver) Years(ctx context.Context) ([]int, error) {
	return r.Repo.Years(ctx)
}

// Area is the resolver for the area field.
func (r *relatedDatasetResolver) Area(ctx context.Context, obj *RelatedDataset) (Area, error) {
	return to[Area](r.Repo.Node(ctx, obj.AreaID))
}

// Type is the resolver for the type field.
func (r *relatedDatasetResolver) Type(ctx context.Context, obj *RelatedDataset) (*RelatedDatasetType, error) {
	return to[*RelatedDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *relatedDatasetItemResolver) Parent(ctx context.Context, obj *RelatedDatasetItem) (*RelatedDataset, error) {
	return to[*RelatedDataset](r.Repo.Node(ctx, obj.ParentID))
}

// GenericDataset returns GenericDatasetResolver implementation.
func (r *Resolver) GenericDataset() GenericDatasetResolver { return &genericDatasetResolver{r} }

// GenericDatasetItem returns GenericDatasetItemResolver implementation.
func (r *Resolver) GenericDatasetItem() GenericDatasetItemResolver {
	return &genericDatasetItemResolver{r}
}

// Municipality returns MunicipalityResolver implementation.
func (r *Resolver) Municipality() MunicipalityResolver { return &municipalityResolver{r} }

// PlateauDataset returns PlateauDatasetResolver implementation.
func (r *Resolver) PlateauDataset() PlateauDatasetResolver { return &plateauDatasetResolver{r} }

// PlateauDatasetItem returns PlateauDatasetItemResolver implementation.
func (r *Resolver) PlateauDatasetItem() PlateauDatasetItemResolver {
	return &plateauDatasetItemResolver{r}
}

// Prefecture returns PrefectureResolver implementation.
func (r *Resolver) Prefecture() PrefectureResolver { return &prefectureResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// RelatedDataset returns RelatedDatasetResolver implementation.
func (r *Resolver) RelatedDataset() RelatedDatasetResolver { return &relatedDatasetResolver{r} }

// RelatedDatasetItem returns RelatedDatasetItemResolver implementation.
func (r *Resolver) RelatedDatasetItem() RelatedDatasetItemResolver {
	return &relatedDatasetItemResolver{r}
}

type genericDatasetResolver struct{ *Resolver }
type genericDatasetItemResolver struct{ *Resolver }
type municipalityResolver struct{ *Resolver }
type plateauDatasetResolver struct{ *Resolver }
type plateauDatasetItemResolver struct{ *Resolver }
type prefectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type relatedDatasetResolver struct{ *Resolver }
type relatedDatasetItemResolver struct{ *Resolver }
