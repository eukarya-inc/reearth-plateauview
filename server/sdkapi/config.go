package sdkapi

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const modelKey = "plateau"

type Config struct {
	CMSBaseURL string
	Project    string
	Model      string
	Token      string
}

func (c *Config) Normalize() {
	if c.Model == "" {
		c.Model = modelKey
	}
}

type DatasetResponse struct {
	Data []DatasetPref `json:"data"`
}

type DatasetPref struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Data  []DatasetCity
}

type DatasetCity struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	FeatureTypes []string `json:"feature_types"`
}

type FilesResponse map[string][]File

type File struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	MaxLOD string `json:"max_lod"`
}

type Items []Item

func (i Items) DatasetResponse() (r *DatasetResponse) {
	r = &DatasetResponse{}
	prefs := []DatasetPref{}
	prefm := map[string]*DatasetPref{}
	for _, i := range i {
		if _, ok := prefm[i.Prefecture]; !ok {
			pd := DatasetPref{
				ID:    i.Prefecture,
				Title: i.Prefecture,
			}
			prefs = append(prefs, pd)
			prefm[i.Prefecture] = lo.ToPtr(prefs[len(prefs)-1])
		}

		d := DatasetCity{
			ID:           i.ID,
			Title:        i.CityName,
			Description:  i.Description,
			FeatureTypes: i.FeatureTypes(),
		}
		pd := prefm[i.Prefecture]
		pd.Data = append(pd.Data, d)
	}

	r.Data = prefs
	return
}

type Item struct {
	ID          string  `json:"id"`
	Prefecture  string  `json:"prefecture"`
	CityName    string  `json:"city_name"`
	CityGML     *Asset  `json:"citygml"`
	Description string  `json:"description_bldg"`
	MaxLOD      *Asset  `json:"max_lod"`
	Bldg        []Asset `json:"bldg"`
	Tran        []Asset `json:"tran"`
	Frn         []Asset `json:"frn"`
	Veg         []Asset `json:"veg"`
}

func (i Item) FeatureTypes() (t []string) {
	if len(i.Bldg) > 0 {
		t = append(t, "bldg")
	}
	if len(i.Tran) > 0 {
		t = append(t, "tran")
	}
	if len(i.Frn) > 0 {
		t = append(t, "frn")
	}
	if len(i.Veg) > 0 {
		t = append(t, "veg")
	}
	return
}

type Asset struct {
	URL string `json:"url"`
}

type MaxLODColumns []MaxLODColumn

type MaxLODMap map[string]map[string]float64

func (mc MaxLODColumns) Map() MaxLODMap {
	m := MaxLODMap{}

	for _, c := range mc {
		max := c.MaxLODAsFloat64()
		if max == 0 {
			continue
		}

		if _, ok := m[c.Type]; !ok {
			m[c.Type] = map[string]float64{}
		}
		t := m[c.Type]
		t[c.Code] = max
	}

	return m
}

func (mm MaxLODMap) ForEachType() map[string]float64 {
	m := map[string]float64{}

	for ty, c := range mm {
		max := 0.0
		for _, lod := range c {
			if lod > max {
				max = lod
			}
		}

		m[ty] = max
	}

	return m
}

func (mm MaxLODMap) Files(citygmlAssetURL string) (r FilesResponse) {
	r = FilesResponse{}
	for ty, m := range mm {
		if _, ok := r[ty]; !ok {
			r[ty] = ([]File)(nil)
		}
		for code, maxlod := range m {
			r[ty] = append(r[ty], File{
				Code:   code,
				URL:    cityGMLURLFromAsset(citygmlAssetURL, ty, code),
				MaxLOD: fmt.Sprintf("%f", maxlod),
			})
		}
	}
	return
}

type MaxLODColumn struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	MaxLOD string `json:"max_lod"`
}

func (m MaxLODColumn) MaxLODAsInt() int {
	n, _ := strconv.Atoi(m.MaxLOD)
	return n
}

func (m MaxLODColumn) MaxLODAsFloat64() float64 {
	n, _ := strconv.ParseFloat(m.MaxLOD, 64)
	return n
}

func cityGMLURLFromAsset(u, ty, code string) string {
	v, err := url.Parse(u)
	if err != nil {
		return ""
	}

	dir := path.Dir(v.Path)
	filename := strings.TrimSuffix(path.Base(v.Path), path.Ext(v.Path))

	fn := fmt.Sprintf("%s_%s_6697_op.gml", code, ty)
	v.Path = path.Join(dir, filename, "udx", ty, fn)
	return v.String()
}
