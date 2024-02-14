package sdkapiv3

import (
	"github.com/hasura/go-graphql-client"
)

type Query struct {
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
