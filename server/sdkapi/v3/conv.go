package sdkapiv3

func (d *Query) ToDatasets() *DatasetsResponse {
	datasets := &DatasetsResponse{}

	for _, prefecture := range d.Areas {
		p := &DatasetPrefectureResponse{
			Title: string(prefecture.Name),
		}

		for _, city := range prefecture.Prefecture.Cities {
			c := &DatasetCityResponse{
				ID:          string(city.ID),
				Title:       string(city.Name),
				Spec:        "", //TODO
				Description: "", //TODO
			}

			for _, dataset := range city.Datasets {
				c.FeatureTypes = append(c.FeatureTypes, string(dataset.TypeCode))
			}

			p.Data = append(p.Data, c)
		}

		datasets.Data = append(datasets.Data, p)
	}

	return datasets
}
