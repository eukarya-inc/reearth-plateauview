package plateauapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/samber/lo"
)

// Prefecture is the resolver for the prefecture field.
func (r *cityResolver) Prefecture(ctx context.Context, obj *City) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// Wards is the resolver for the wards field.
func (r *cityResolver) Wards(ctx context.Context, obj *City) ([]*Ward, error) {
	areas, err := r.Repo.Areas(ctx, AreaQuery{
		ParentCode: lo.ToPtr(obj.Code),
	})
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(areas, func(a Area, _ int) (*Ward, bool) {
		if m, ok := a.(*Ward); ok {
			return m, ok
		}
		return nil, false
	}), nil
}

// Datasets is the resolver for the datasets field.
func (r *cityResolver) Datasets(ctx context.Context, obj *City, input DatasetForAreaQuery) ([]Dataset, error) {
	return r.Repo.Datasets(ctx, DatasetQuery{
		AreaCodes:    []AreaCode{obj.Code},
		ExcludeTypes: input.ExcludeTypes,
		IncludeTypes: input.IncludeTypes,
		SearchTokens: input.SearchTokens,
	})
}

// Prefecture is the resolver for the prefecture field.
func (r *genericDatasetResolver) Prefecture(ctx context.Context, obj *GenericDataset) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// City is the resolver for the city field.
func (r *genericDatasetResolver) City(ctx context.Context, obj *GenericDataset) (*City, error) {
	if obj.CityID == nil {
		return nil, nil
	}
	return to[*City](r.Repo.Node(ctx, *obj.CityID))
}

// Ward is the resolver for the ward field.
func (r *genericDatasetResolver) Ward(ctx context.Context, obj *GenericDataset) (*Ward, error) {
	if obj.WardID == nil {
		return nil, nil
	}
	return to[*Ward](r.Repo.Node(ctx, *obj.WardID))
}

// Type is the resolver for the type field.
func (r *genericDatasetResolver) Type(ctx context.Context, obj *GenericDataset) (*GenericDatasetType, error) {
	return to[*GenericDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *genericDatasetItemResolver) Parent(ctx context.Context, obj *GenericDatasetItem) (*GenericDataset, error) {
	return to[*GenericDataset](r.Repo.Node(ctx, obj.ParentID))
}

// Prefecture is the resolver for the prefecture field.
func (r *plateauDatasetResolver) Prefecture(ctx context.Context, obj *PlateauDataset) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// City is the resolver for the city field.
func (r *plateauDatasetResolver) City(ctx context.Context, obj *PlateauDataset) (*City, error) {
	if obj.CityID == nil {
		return nil, nil
	}
	return to[*City](r.Repo.Node(ctx, *obj.CityID))
}

// Ward is the resolver for the ward field.
func (r *plateauDatasetResolver) Ward(ctx context.Context, obj *PlateauDataset) (*Ward, error) {
	if obj.WardID == nil {
		return nil, nil
	}
	return to[*Ward](r.Repo.Node(ctx, *obj.WardID))
}

// Type is the resolver for the type field.
func (r *plateauDatasetResolver) Type(ctx context.Context, obj *PlateauDataset) (*PlateauDatasetType, error) {
	return to[*PlateauDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *plateauDatasetItemResolver) Parent(ctx context.Context, obj *PlateauDatasetItem) (*PlateauDataset, error) {
	return to[*PlateauDataset](r.Repo.Node(ctx, obj.ParentID))
}

// PlateauSpec is the resolver for the plateauSpec field.
func (r *plateauDatasetTypeResolver) PlateauSpec(ctx context.Context, obj *PlateauDatasetType) (*PlateauSpec, error) {
	return to[*PlateauSpec](r.Repo.Node(ctx, obj.PlateauSpecID))
}

// Prefecture is the resolver for the prefecture field.
func (r *plateauFloodingDatasetResolver) Prefecture(ctx context.Context, obj *PlateauFloodingDataset) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// City is the resolver for the city field.
func (r *plateauFloodingDatasetResolver) City(ctx context.Context, obj *PlateauFloodingDataset) (*City, error) {
	if obj.CityID == nil {
		return nil, nil
	}
	return to[*City](r.Repo.Node(ctx, *obj.CityID))
}

// Ward is the resolver for the ward field.
func (r *plateauFloodingDatasetResolver) Ward(ctx context.Context, obj *PlateauFloodingDataset) (*Ward, error) {
	if obj.WardID == nil {
		return nil, nil
	}
	return to[*Ward](r.Repo.Node(ctx, *obj.WardID))
}

// Type is the resolver for the type field.
func (r *plateauFloodingDatasetResolver) Type(ctx context.Context, obj *PlateauFloodingDataset) (*PlateauDatasetType, error) {
	return to[*PlateauDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *plateauFloodingDatasetItemResolver) Parent(ctx context.Context, obj *PlateauFloodingDatasetItem) (*PlateauDataset, error) {
	return to[*PlateauDataset](r.Repo.Node(ctx, obj.ParentID))
}

// Cities is the resolver for the cities field.
func (r *prefectureResolver) Cities(ctx context.Context, obj *Prefecture) ([]*City, error) {
	areas, err := r.Repo.Areas(ctx, AreaQuery{
		ParentCode: lo.ToPtr(obj.Code),
	})
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(areas, func(a Area, _ int) (*City, bool) {
		if m, ok := a.(*City); ok {
			return m, ok
		}
		return nil, false
	}), nil
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

// Prefecture is the resolver for the prefecture field.
func (r *relatedDatasetResolver) Prefecture(ctx context.Context, obj *RelatedDataset) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// City is the resolver for the city field.
func (r *relatedDatasetResolver) City(ctx context.Context, obj *RelatedDataset) (*City, error) {
	if obj.CityID == nil {
		return nil, nil
	}
	return to[*City](r.Repo.Node(ctx, *obj.CityID))
}

// Ward is the resolver for the ward field.
func (r *relatedDatasetResolver) Ward(ctx context.Context, obj *RelatedDataset) (*Ward, error) {
	if obj.WardID == nil {
		return nil, nil
	}
	return to[*Ward](r.Repo.Node(ctx, *obj.WardID))
}

// Type is the resolver for the type field.
func (r *relatedDatasetResolver) Type(ctx context.Context, obj *RelatedDataset) (*RelatedDatasetType, error) {
	return to[*RelatedDatasetType](r.Repo.Node(ctx, obj.TypeID))
}

// Parent is the resolver for the parent field.
func (r *relatedDatasetItemResolver) Parent(ctx context.Context, obj *RelatedDatasetItem) (*RelatedDataset, error) {
	return to[*RelatedDataset](r.Repo.Node(ctx, obj.ParentID))
}

// Prefecture is the resolver for the prefecture field.
func (r *wardResolver) Prefecture(ctx context.Context, obj *Ward) (*Prefecture, error) {
	return to[*Prefecture](r.Repo.Node(ctx, obj.PrefectureID))
}

// City is the resolver for the city field.
func (r *wardResolver) City(ctx context.Context, obj *Ward) (*City, error) {
	return to[*City](r.Repo.Node(ctx, obj.CityID))
}

// Datasets is the resolver for the datasets field.
func (r *wardResolver) Datasets(ctx context.Context, obj *Ward, input DatasetForAreaQuery) ([]Dataset, error) {
	return r.Repo.Datasets(ctx, DatasetQuery{
		AreaCodes:    []AreaCode{obj.Code},
		ExcludeTypes: input.ExcludeTypes,
		IncludeTypes: input.IncludeTypes,
		SearchTokens: input.SearchTokens,
	})
}

// City returns CityResolver implementation.
func (r *Resolver) City() CityResolver { return &cityResolver{r} }

// GenericDataset returns GenericDatasetResolver implementation.
func (r *Resolver) GenericDataset() GenericDatasetResolver { return &genericDatasetResolver{r} }

// GenericDatasetItem returns GenericDatasetItemResolver implementation.
func (r *Resolver) GenericDatasetItem() GenericDatasetItemResolver {
	return &genericDatasetItemResolver{r}
}

// PlateauDataset returns PlateauDatasetResolver implementation.
func (r *Resolver) PlateauDataset() PlateauDatasetResolver { return &plateauDatasetResolver{r} }

// PlateauDatasetItem returns PlateauDatasetItemResolver implementation.
func (r *Resolver) PlateauDatasetItem() PlateauDatasetItemResolver {
	return &plateauDatasetItemResolver{r}
}

// PlateauDatasetType returns PlateauDatasetTypeResolver implementation.
func (r *Resolver) PlateauDatasetType() PlateauDatasetTypeResolver {
	return &plateauDatasetTypeResolver{r}
}

// PlateauFloodingDataset returns PlateauFloodingDatasetResolver implementation.
func (r *Resolver) PlateauFloodingDataset() PlateauFloodingDatasetResolver {
	return &plateauFloodingDatasetResolver{r}
}

// PlateauFloodingDatasetItem returns PlateauFloodingDatasetItemResolver implementation.
func (r *Resolver) PlateauFloodingDatasetItem() PlateauFloodingDatasetItemResolver {
	return &plateauFloodingDatasetItemResolver{r}
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

// Ward returns WardResolver implementation.
func (r *Resolver) Ward() WardResolver { return &wardResolver{r} }

type cityResolver struct{ *Resolver }
type genericDatasetResolver struct{ *Resolver }
type genericDatasetItemResolver struct{ *Resolver }
type plateauDatasetResolver struct{ *Resolver }
type plateauDatasetItemResolver struct{ *Resolver }
type plateauDatasetTypeResolver struct{ *Resolver }
type plateauFloodingDatasetResolver struct{ *Resolver }
type plateauFloodingDatasetItemResolver struct{ *Resolver }
type prefectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type relatedDatasetResolver struct{ *Resolver }
type relatedDatasetItemResolver struct{ *Resolver }
type wardResolver struct{ *Resolver }
