package sdkapiv3

import "github.com/hasura/go-graphql-client"

const bldg = "bldg"

func (d *DatasetsQuery) ToDatasets() *DatasetsResponse {
	datasets := &DatasetsResponse{}

	for _, prefecture := range d.Areas {
		p := &DatasetPrefectureResponse{
			Title: string(prefecture.Name),
		}

		for _, city := range prefecture.Prefecture.Cities {
			if city.Citygml == nil || len(city.Datasets) == 0 {
				continue
			}

			c := &DatasetCityResponse{
				ID:           string(city.Code),
				Title:        string(city.Name),
				FeatureTypes: toStrings(city.Citygml.FeatureTypes),
				Spec:         string(city.Citygml.PlateauSpecMinor.Version),
			}

			for _, dataset := range city.Datasets {
				if dataset.TypeCode == bldg {
					c.Description = string(dataset.Description)
					break
				}
			}

			p.Data = append(p.Data, c)
		}

		if len(p.Data) > 0 {
			datasets.Data = append(datasets.Data, p)
		}
	}

	return datasets
}

func toStrings(s []graphql.String) []string {
	var r []string
	for _, v := range s {
		r = append(r, string(v))
	}
	return r
}
