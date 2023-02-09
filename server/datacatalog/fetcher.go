package datacatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

const ModelPlateau = "plateau"
const ModelUsecase = "usecase"
const ModelDataset = "dataset"

type Config struct {
	CMSBase    string
	CMSProject string
}

type Fetcher struct {
	c    *http.Client
	base *url.URL
}

func NewFetcher(c *http.Client, config Config) (*Fetcher, error) {
	u, err := url.Parse(config.CMSBase)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "api", "p", config.CMSProject)

	if c == nil {
		c = http.DefaultClient
	}

	return &Fetcher{
		c:    c,
		base: u,
	}, nil
}

func (f *Fetcher) Do(ctx context.Context) (DataCatalogItems, error) {
	resultPlateau := make(chan []DataCatalogItem)
	resultUsecase := make(chan []DataCatalogItem)
	resultDataset := make(chan []DataCatalogItem)
	errPlateau := make(chan error)
	errUsecase := make(chan error)
	errDataset := make(chan error)

	go func() {
		r, err := f.all(ctx, ModelPlateau)
		errPlateau <- err
		resultPlateau <- r
	}()

	go func() {
		r, err := f.all(ctx, ModelUsecase)
		errUsecase <- err
		resultUsecase <- r
	}()

	go func() {
		r, err := f.all(ctx, ModelDataset)
		errDataset <- err
		resultDataset <- r
	}()

	if err := <-errPlateau; err != nil {
		return DataCatalogItems{}, err
	}

	if err := <-errUsecase; err != nil {
		return DataCatalogItems{}, err
	}

	if err := <-errDataset; err != nil {
		return DataCatalogItems{}, err
	}

	return DataCatalogItems{
		Plateau: <-resultPlateau,
		Usecase: <-resultUsecase,
		Dataset: <-resultDataset,
	}, nil
}

func (f *Fetcher) all(ctx context.Context, model string) (resp []DataCatalogItem, err error) {
	for p := 1; ; p++ {
		r, err := f.get(ctx, model, p, 0)
		if err != nil {
			return nil, err
		}

		if !r.HasNext() {
			break
		}

		resp = append(resp, r.DataCatalogs()...)
	}
	return
}

func (f *Fetcher) get(ctx context.Context, model string, page, perPage int) (r Response, err error) {
	if perPage == 0 {
		perPage = 100
	}

	req, err := http.NewRequestWithContext(ctx, "GET", f.url(model, page, perPage), nil)
	if err != nil {
		return
	}

	res, err := f.c.Do(req)
	if err != nil {
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %d", res.StatusCode)
		return
	}

	err = json.NewDecoder(res.Request.Body).Decode(&r)
	r.Page = page
	r.PerPage = perPage
	return
}

func (f *Fetcher) url(model string, page, perPage int) string {
	u := util.CloneRef(f.base)
	u.Path = path.Join(u.Path, model)
	u.RawQuery = url.Values{
		"page":    []string{strconv.Itoa(page)},
		"perPage": []string{strconv.Itoa(perPage)},
	}.Encode()
	return u.String()
}

type Response struct {
	Model      string          `json:"-"`
	Results    json.RawMessage `json:"results"`
	Plateau    []plateauItem   `json:"-"`
	Usecase    []usecaseItem   `json:"-"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	TotalCount int             `json:"total_count"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	if r.Model == ModelPlateau {
		if err := json.Unmarshal(r.Results, &r.Plateau); err != nil {
			return err
		}
	}
	if r.Model == ModelUsecase || r.Model == ModelDataset {
		if err := json.Unmarshal(r.Results, &r.Usecase); err != nil {
			return err
		}
	}
	return nil
}

func (r Response) HasNext() bool {
	return r.TotalCount/r.PerPage >= r.Page
}

func (r Response) DataCatalogs() []DataCatalogItem {
	if r.Plateau != nil {
		return lo.FlatMap(r.Plateau, func(i plateauItem, _ int) []DataCatalogItem {
			return i.DataCatalogs()
		})
	}
	if r.Usecase != nil {
		return lo.FlatMap(r.Usecase, func(i usecaseItem, _ int) []DataCatalogItem {
			return i.DataCatalogs()
		})
	}
	return nil
}
