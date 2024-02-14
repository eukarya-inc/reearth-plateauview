package sdkapiv3

import (
	"github.com/hasura/go-graphql-client"
)

type Query struct {
	Areas        []*QueryArea `graphql:"areas(input: {areaTypes: PREFECTURE})"`
	PlateauSpecs []*QueryPlateauSpec
}

type QueryArea struct {
	ID         graphql.String
	Name       graphql.String
	Prefecture *QueryPrefecture `graphql:"... on Prefecture"`
}

type QueryPrefecture struct {
	Cities []*QueryCity
}

type QueryCity struct {
	ID       graphql.String
	Name     graphql.String
	Datasets []*QueryCityDataset `graphql:"datasets(input: {includeTypes: [\"plateau\"]})"`
}

type QueryCityDataset struct {
	ID          graphql.String
	Name        graphql.String
	TypeCode    graphql.String
	Description graphql.String
}

type QueryPlateauSpec struct {
	MajorVersion  graphql.Int
	MinorVersions []*QueryPlateauSpecMinorVersion
}

type QueryPlateauSpecMinorVersion struct {
	Version graphql.String
}
