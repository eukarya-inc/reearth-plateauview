//go:generate go run github.com/99designs/gqlgen generate

package plateauapi

import "github.com/99designs/gqlgen/graphql"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func NewSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{Resolvers: &Resolver{}})
}

func (r *Resolver) Country() CountryResolver           { return &areaResolver[*Country]{r} }
func (r *Resolver) Municipality() MunicipalityResolver { return &areaResolver[*Municipality]{r} }
func (r *Resolver) Prefecture() PrefectureResolver     { return &areaResolver[*Prefecture]{r} }

type areaResolver[T Area] struct{ *Resolver }

func (r *Resolver) GenericDataset() GenericDatasetResolver {
	return &datasetResolver[*GenericDataset]{r}

}
func (r *Resolver) PlateauDataset() PlateauDatasetResolver {
	return &datasetResolver[*PlateauDataset]{r}
}

func (r *Resolver) PlateauAuxiliaryDataset() PlateauAuxiliaryDatasetResolver {
	return &datasetResolver[*PlateauAuxiliaryDataset]{r}
}

type datasetResolver[T Dataset] struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
