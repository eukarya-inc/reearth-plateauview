package sdkapiv3

import (
	"sort"

	"golang.org/x/exp/maps"
)

const bldg = "bldg"

func (d *DatasetsQuery) ToDatasets() *DatasetsResponse {
	datasets := &DatasetsResponse{}

	for _, prefecture := range d.Areas {
		p := &DatasetPrefectureResponse{
			Title: string(prefecture.Name),
		}

		for _, city := range prefecture.Prefecture.Cities {
			c := &DatasetCityResponse{
				ID:    string(city.Code),
				Title: string(city.Name),
			}

			ft := map[string]struct{}{}

			for _, dataset := range city.Datasets {
				if dataset.TypeCode == bldg {
					c.Description = string(dataset.Description)
					c.Spec = string(dataset.PlateauDataset.PlateauSpecMinor.Version)
				}
				ft[string(dataset.TypeCode)] = struct{}{}
			}

			c.FeatureTypes = maps.Keys(ft)
			sort.Strings(c.FeatureTypes)

			p.Data = append(p.Data, c)
		}

		datasets.Data = append(datasets.Data, p)
	}

	return datasets
}

func (d *DatasetFilesQuery) ToDatasetFiles() *map[string][]DatasetFilesResponse {
	files := make(map[string][]DatasetFilesResponse)

	for _, item := range d.Area.City.Citygml.Items {
		files[string(item.TypeCode)] = append(files[string(item.TypeCode)], DatasetFilesResponse{
			Code:   string(item.MeshCode),
			URL:    string(item.Url),
			MaxLod: int(item.MaxLod),
		})
	}

	return &files
}
