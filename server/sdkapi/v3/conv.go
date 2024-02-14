package sdkapiv3

const bldg = "bldg"

func (d *Query) ToDatasets() *DatasetsResponse {
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
