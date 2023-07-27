package datacatalog

import (
	"context"
	"errors"
	"net/url"
	"path"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauv2"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

const timeoutSecond int64 = 20
const ModelPlateau = "plateau"
const ModelUsecase = "usecase"
const ModelDataset = "dataset"

type Fetcher struct {
	cmsp *cms.PublicAPIClient[plateauv2.CMSItem]
	cmsu *cms.PublicAPIClient[UsecaseItem]
	base *url.URL
}

type FetcherDoOptions struct {
	Subproject string
	CityName   string
}

func NewFetcher(cmsbase string) (*Fetcher, error) {
	u, err := url.Parse(cmsbase)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "api", "p")

	cmsp, err := cms.NewPublicAPIClient[any](nil, cmsbase)
	if err != nil {
		return nil, err
	}

	cmsp = cmsp.WithTimeout(time.Duration(timeoutSecond) * time.Second)

	return &Fetcher{
		cmsp: cms.ChangePublicAPIClientType[any, plateauv2.CMSItem](cmsp),
		cmsu: cms.ChangePublicAPIClientType[any, UsecaseItem](cmsp),
		base: u,
	}, nil
}

func (f *Fetcher) Clone() *Fetcher {
	if f == nil {
		return nil
	}

	return &Fetcher{
		cmsp: f.cmsp.Clone(),
		cmsu: f.cmsu.Clone(),
		base: util.CloneRef(f.base),
	}
}

func (f *Fetcher) Do(ctx context.Context, project string, opts FetcherDoOptions) (ResponseAll, error) {
	f1, f2, f3, f4, f5 := f.Clone(), f.Clone(), f.Clone(), f.Clone(), f.Clone()

	res1 := lo.Async2(func() ([]plateauv2.CMSItem, error) {
		return f1.plateau(ctx, project, ModelPlateau)
	})
	res2 := lo.Async2(func() ([]UsecaseItem, error) {
		return f2.usecase(ctx, project, ModelUsecase)
	})
	res3 := lo.Async2(func() ([]UsecaseItem, error) {
		return f3.usecase(ctx, project, ModelDataset)
	})
	res4 := lo.Async2(func() ([]plateauv2.CMSItem, error) {
		if opts.CityName == "" || opts.Subproject == "" {
			return nil, nil
		}
		return f4.plateau(ctx, opts.Subproject, ModelPlateau)
	})
	res5 := lo.Async2(func() ([]UsecaseItem, error) {
		if opts.CityName == "" || opts.Subproject == "" {
			return nil, nil
		}
		return f5.usecase(ctx, opts.Subproject, ModelDataset)
	})

	notFound := 0
	r := ResponseAll{}

	if res := <-res1; res.B != nil {
		if errors.Is(res.B, cms.ErrNotFound) {
			notFound++
		} else {
			return ResponseAll{}, res.B
		}
	} else {
		r.Plateau = append(r.Plateau, res.A...)
	}

	if res := <-res2; res.B != nil {
		if errors.Is(res.B, cms.ErrNotFound) {
			notFound++
		} else {
			return ResponseAll{}, res.B
		}
	} else {
		r.Usecase = append(r.Usecase, res.A...)
	}

	if res := <-res3; res.B != nil {
		if errors.Is(res.B, cms.ErrNotFound) {
			notFound++
		} else {
			return ResponseAll{}, res.B
		}
	} else {
		r.Usecase = append(r.Usecase, res.A...)
	}

	if res := <-res4; res.B != nil {
		if errors.Is(res.B, cms.ErrNotFound) {
			notFound++
		} else {
			return ResponseAll{}, res.B
		}
	} else {
		r.Plateau = append(r.Plateau, filterPlateau(res.A, opts.CityName)...)
	}

	if res := <-res5; res.B != nil {
		if errors.Is(res.B, cms.ErrNotFound) {
			notFound++
		} else {
			return ResponseAll{}, res.B
		}
	} else {
		r.Usecase = append(r.Usecase, filterUsecase(res.A, opts.CityName)...)
	}

	if notFound == 3 {
		return r, rerror.ErrNotFound
	}
	return r, nil
}

func (f *Fetcher) plateau(ctx context.Context, project, model string) (resp []plateauv2.CMSItem, err error) {
	r, err := f.cmsp.GetAllItemsInParallel(ctx, project, model, 10)
	if err != nil {
		return
	}
	return r, nil
}

func (f *Fetcher) usecase(ctx context.Context, project, model string) (resp []UsecaseItem, err error) {
	r, err := f.cmsu.GetAllItemsInParallel(ctx, project, model, 10)
	if err != nil {
		return
	}

	if model == ModelUsecase {
		for i := range r {
			r[i].Type = "ユースケース"
		}
	}

	return r, nil
}

func filterPlateau(data []plateauv2.CMSItem, cityName string) []plateauv2.CMSItem {
	if cityName == "" {
		return nil
	}

	return lo.Filter(data, func(v plateauv2.CMSItem, _ int) bool {
		return v.CityName == cityName
	})
}

func filterUsecase(data []UsecaseItem, cityName string) []UsecaseItem {
	if cityName == "" {
		return nil
	}

	return lo.Filter(data, func(v UsecaseItem, _ int) bool {
		return v.CityName == cityName
	})
}
