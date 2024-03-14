package datacatalog

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/spkg/bom"
)

type CityGMLFilesResponse struct {
	CityCode         string       `json:"cityCode"`
	CityName         string       `json:"cityName"`
	Year             int          `json:"year"`
	RegistrationYear int          `json:"registrationYear"`
	Spec             string       `json:"spec"`
	URL              string       `json:"url"`
	Files            CityGMLFiles `json:"files"`
}

type CityGMLFiles = map[string][]CityGMLFile

type CityGMLFile struct {
	MeshCode string `json:"code"`
	MaxLOD   int    `json:"maxLod"`
	URL      string `json:"url"`
}

func fetchCityGMLFiles(ctx context.Context, r plateauapi.Repo, id string) (*CityGMLFilesResponse, error) {
	n, err := r.Node(ctx, plateauapi.CityGMLDatasetIDFrom(plateauapi.AreaCode(id)))
	if err != nil {
		return nil, err
	}

	citygml, ok := n.(*plateauapi.CityGMLDataset)
	if !ok || citygml == nil || citygml.URL == "" || citygml.PlateauSpecMinorID == "" {
		return nil, nil
	}

	n, err = r.Node(ctx, citygml.PlateauSpecMinorID)
	if err != nil {
		return nil, err
	}

	spec, ok := n.(*plateauapi.PlateauSpecMinor)
	if !ok || spec == nil {
		return nil, nil
	}

	n, err = r.Node(ctx, citygml.CityID)
	if err != nil {
		return nil, err
	}

	city, ok := n.(*plateauapi.City)
	if !ok || city == nil {
		return nil, nil
	}

	admin, ok := citygml.Admin.(map[string]any)
	if !ok || admin == nil {
		return nil, nil
	}

	url := admin["maxlod"].(string)
	if url == "" {
		return nil, nil
	}

	data, err := fetchCSV(ctx, url)
	if err != nil {
		return nil, err
	}

	files := csvToCityGMLFilesResponse(data, citygml.URL)
	return &CityGMLFilesResponse{
		CityCode:         string(citygml.CityCode),
		CityName:         city.Name,
		Year:             citygml.Year,
		RegistrationYear: citygml.RegistrationYear,
		Spec:             spec.Version,
		URL:              citygml.URL,
		Files:            files,
	}, nil
}

func fetchCSV(ctx context.Context, url string) (records [][]string, _ error) {
	res, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(res)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	c := csv.NewReader(bom.NewReader(resp.Body))
	for {
		record, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read csv: %w", err)
		}
		records = append(records, record)
	}

	return
}

func csvToCityGMLFilesResponse(data [][]string, citygmlURL string) CityGMLFiles {
	res := make(CityGMLFiles)

	for _, record := range data {
		if len(record) < 3 || record[0] == "" {
			continue
		}

		if !isNumeric(rune(record[0][0])) {
			// it's a header
			continue
		}

		// code,type,maxLod(,path)
		meshCode := record[0]
		featureType := record[1]
		maxlod, _ := strconv.Atoi(record[2])
		url := ""
		if len(record) > 3 {
			url = citygmlItemURLFrom(citygmlURL, record[3], featureType)
			if url == "" {
				continue
			}
		}

		item := CityGMLFile{
			MeshCode: meshCode,
			MaxLOD:   maxlod,
			URL:      url,
		}

		if _, ok := res[featureType]; !ok {
			res[featureType] = make([]CityGMLFile, 0)
		}

		res[featureType] = append(res[featureType], item)
	}

	return res
}

func citygmlItemURLFrom(base, p, typeCode string) string {
	b := path.Base(base)
	base = strings.TrimSuffix(base, b)
	u, _ := url.JoinPath(base, nameWithoutExt(b), "udx", typeCode, p)
	return u
}

func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func nameWithoutExt(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}
