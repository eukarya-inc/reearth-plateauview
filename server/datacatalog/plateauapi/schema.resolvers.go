package plateauapi

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

// Datasets is the resolver for the datasets field.
func (r *countryResolver) Datasets(ctx context.Context, obj *Country, input DatasetForAreaQuery) ([]Dataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Area is the resolver for the area field.
func (r *genericDatasetResolver) Area(ctx context.Context, obj *GenericDataset) (Area, error) {
	panic(fmt.Errorf("not implemented"))
}

// Type is the resolver for the type field.
func (r *genericDatasetResolver) Type(ctx context.Context, obj *GenericDataset) (*GenericDatasetType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Parent is the resolver for the parent field.
func (r *genericDatasetDatumResolver) Parent(ctx context.Context, obj *GenericDatasetDatum) (*GenericDataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Datasets is the resolver for the datasets field.
func (r *municipalityResolver) Datasets(ctx context.Context, obj *Municipality, input DatasetForAreaQuery) ([]Dataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Area is the resolver for the area field.
func (r *plateauAuxiliaryDatasetResolver) Area(ctx context.Context, obj *PlateauAuxiliaryDataset) (Area, error) {
	panic(fmt.Errorf("not implemented"))
}

// Type is the resolver for the type field.
func (r *plateauAuxiliaryDatasetResolver) Type(ctx context.Context, obj *PlateauAuxiliaryDataset) (*PlateauAuxiliaryDatasetType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Parent is the resolver for the parent field.
func (r *plateauAuxiliaryDatasetDatumResolver) Parent(ctx context.Context, obj *PlateauAuxiliaryDatasetDatum) (*PlateauAuxiliaryDataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Area is the resolver for the area field.
func (r *plateauDatasetResolver) Area(ctx context.Context, obj *PlateauDataset) (Area, error) {
	panic(fmt.Errorf("not implemented"))
}

// Type is the resolver for the type field.
func (r *plateauDatasetResolver) Type(ctx context.Context, obj *PlateauDataset) (*PlateauDatasetType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Parent is the resolver for the parent field.
func (r *plateauDatasetDatumResolver) Parent(ctx context.Context, obj *PlateauDatasetDatum) (*PlateauDataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Datasets is the resolver for the datasets field.
func (r *prefectureResolver) Datasets(ctx context.Context, obj *Prefecture, input DatasetForAreaQuery) ([]Dataset, error) {
	panic(fmt.Errorf("not implemented"))
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (Node, error) {
	panic("implement me")
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]Node, error) {
	panic("implement me")
}

// Area is the resolver for the area field.
func (r *queryResolver) Area(ctx context.Context, code string) (Area, error) {
	panic("implement me")
}

// Areas is the resolver for the areas field.
func (r *queryResolver) Areas(ctx context.Context, input AreaQuery) ([]Area, error) {
	panic("implement me")
}

// DatasetTypes is the resolver for the datasetTypes field.
func (r *queryResolver) DatasetTypes(ctx context.Context, input DatasetTypeQuery) ([]DatasetType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Datasets is the resolver for the datasets field.
func (r *queryResolver) Datasets(ctx context.Context, input DatasetQuery) ([]Dataset, error) {
	panic("implement me")
}

// PlateauVersions is the resolver for the plateauVersions field.
func (r *queryResolver) PlateauVersions(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Years is the resolver for the years field.
func (r *queryResolver) Years(ctx context.Context) ([]int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Country returns CountryResolver implementation.
func (r *Resolver) Country() CountryResolver { return &countryResolver{r} }

// GenericDataset returns GenericDatasetResolver implementation.
func (r *Resolver) GenericDataset() GenericDatasetResolver { return &genericDatasetResolver{r} }

// GenericDatasetDatum returns GenericDatasetDatumResolver implementation.
func (r *Resolver) GenericDatasetDatum() GenericDatasetDatumResolver {
	return &genericDatasetDatumResolver{r}
}

// Municipality returns MunicipalityResolver implementation.
func (r *Resolver) Municipality() MunicipalityResolver { return &municipalityResolver{r} }

// PlateauAuxiliaryDataset returns PlateauAuxiliaryDatasetResolver implementation.
func (r *Resolver) PlateauAuxiliaryDataset() PlateauAuxiliaryDatasetResolver {
	return &plateauAuxiliaryDatasetResolver{r}
}

// PlateauAuxiliaryDatasetDatum returns PlateauAuxiliaryDatasetDatumResolver implementation.
func (r *Resolver) PlateauAuxiliaryDatasetDatum() PlateauAuxiliaryDatasetDatumResolver {
	return &plateauAuxiliaryDatasetDatumResolver{r}
}

// PlateauDataset returns PlateauDatasetResolver implementation.
func (r *Resolver) PlateauDataset() PlateauDatasetResolver { return &plateauDatasetResolver{r} }

// PlateauDatasetDatum returns PlateauDatasetDatumResolver implementation.
func (r *Resolver) PlateauDatasetDatum() PlateauDatasetDatumResolver {
	return &plateauDatasetDatumResolver{r}
}

// Prefecture returns PrefectureResolver implementation.
func (r *Resolver) Prefecture() PrefectureResolver { return &prefectureResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type countryResolver struct{ *Resolver }
type genericDatasetResolver struct{ *Resolver }
type genericDatasetDatumResolver struct{ *Resolver }
type municipalityResolver struct{ *Resolver }
type plateauAuxiliaryDatasetResolver struct{ *Resolver }
type plateauAuxiliaryDatasetDatumResolver struct{ *Resolver }
type plateauDatasetResolver struct{ *Resolver }
type plateauDatasetDatumResolver struct{ *Resolver }
type prefectureResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
