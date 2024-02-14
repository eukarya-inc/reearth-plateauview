package sdkapiv3

import (
	"github.com/hasura/go-graphql-client"
)

type Query struct {
	Areas        []QueryArea `graphql:"areas(input: {areaTypes: PREFECTURE})"`
	PlateauSpecs []QueryPlateauSpec
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
	Datasets []QueryCityDataset `graphql:"datasets(input: {includeTypes: [\"plateau\"]})"`
}

type QueryCityDataset struct {
	ID          graphql.String
	Name        graphql.String
	TypeCode    graphql.String
	Description graphql.String
}

type QueryPlateauSpec struct {
	MajorVersion  graphql.Int
	MinorVersions []QueryPlateauSpecMinorVersion
}

type QueryPlateauSpecMinorVersion struct {
	Version graphql.String
}

type Datasets struct {
	Data []DatasetPrefecture `json:"data"`
}

type DatasetPrefecture struct {
	Title string        `json:"title"`
	Data  []DatasetCity `json:"data"`
}

type DatasetCity struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Spec         string   `json:"spec"`
	Description  string   `json:"description"`
	FeatureTypes []string `json:"featureTypes"`
}

func (d *Query) ToDatasets() Datasets {
	var datasets Datasets

	for _, prefecture := range d.Areas {
		var p DatasetPrefecture

		p.Title = string(prefecture.Name)

		for _, city := range prefecture.Prefecture.Cities {
			var c DatasetCity

			c.ID = string(city.ID)
			c.Title = string(city.Name)
			c.Spec = ""        //TODO
			c.Description = "" //TODO

			for _, dataset := range city.Datasets {
				c.FeatureTypes = append(c.FeatureTypes, string(dataset.TypeCode))
			}

			p.Data = append(p.Data, c)
		}

		datasets.Data = append(datasets.Data, p)
	}

	return datasets
}
