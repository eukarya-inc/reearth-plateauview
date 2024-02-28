package sdkapiv3

const bldg = "bldg"

func (d *DatasetsQuery) ToDatasets() *DatasetsResponse {
	datasets := &DatasetsResponse{}

	for _, prefecture := range d.Areas {
		p := &DatasetPrefectureResponse{
			Title: string(prefecture.Name),
		}

		for _, city := range prefecture.Prefecture.Cities {
			c := &DatasetCityResponse{
				ID:    string(city.ID),
				Title: string(city.Name),
			}

			for _, dataset := range city.Datasets {
				if dataset.TypeCode == bldg {
					c.Description = string(dataset.Description)
					c.Spec = string(dataset.PlateauDataset.PlateauSpecMinor.Version)
				}
				c.FeatureTypes = append(c.FeatureTypes, string(dataset.TypeCode))
			}

			p.Data = append(p.Data, c)
		}

		datasets.Data = append(datasets.Data, p)
	}

	return datasets
}

func (d *DatasetFilesQuery) ToDatasetFiles() *map[string][]DatasetFilesResponse {
	files := make(map[string][]DatasetFilesResponse)

	for _, item := range d.City.City.Citygml.Items {
		files[string(item.TypeCode)] = append(files[string(item.TypeCode)], DatasetFilesResponse{
			Code:   string(item.MeshCode),
			URL:    string(item.Url),
			MaxLod: int(item.MaxLod),
		})
	}

	return &files
}
