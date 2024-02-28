package sdkapiv3

type DatasetsResponse struct {
	Data []*DatasetPrefectureResponse `json:"data"`
}

type DatasetPrefectureResponse struct {
	Title string                 `json:"title"`
	Data  []*DatasetCityResponse `json:"data"`
}

type DatasetCityResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Spec         string   `json:"spec"`
	Description  string   `json:"description"`
	FeatureTypes []string `json:"featureTypes"`
	URL          string   `json:"url"`
}

type DatasetFilesResponse struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	MaxLod int    `json:"maxLod"`
}
