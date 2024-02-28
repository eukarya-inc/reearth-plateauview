package sdkapiv3

import (
	"github.com/hasura/go-graphql-client"
)

type DatasetsQuery struct {
	Areas []QueryArea `graphql:"areas(input: {areaTypes: PREFECTURE})"`
}

type QueryArea struct {
	ID         graphql.String
	Name       graphql.String
	Prefecture QueryPrefecture `graphql:"... on Prefecture"`
}

type QueryPrefecture struct {
	Cities []QueryCity
}

type QueryCity struct {
	ID       graphql.String
	Name     graphql.String
	Code     graphql.String
	Datasets []QueryDataset `graphql:"datasets(input: {includeTypes: [\"plateau\"]})"`
}

type QueryDataset struct {
	ID             graphql.String
	Name           graphql.String
	TypeCode       graphql.String
	Description    graphql.String
	PlateauDataset QueryPlateauDataset `graphql:"... on PlateauDataset"`
}

type QueryPlateauDataset struct {
	PlateauSpecMinor QueryPlateauSpecMinor
}

type QueryPlateauSpecMinor struct {
	Version graphql.String
}

type DatasetFilesQuery struct {
	Area QueryFilesArea `graphql:"area(code: $code)"`
}

type QueryFilesArea struct {
	City QueryFilesCity `graphql:"... on City"`
}

type QueryFilesCity struct {
	ID      graphql.String
	Citygml QueryFilesCityGML
}

type QueryFilesCityGML struct {
	Items []QueryFilesCityGMLItem
}

type QueryFilesCityGMLItem struct {
	Url      graphql.String
	TypeCode graphql.String
	MeshCode graphql.String
	MaxLod   graphql.Int
}

type AreaCode string

func (AreaCode) GetGraphQLType() string { return "AreaCode" }
