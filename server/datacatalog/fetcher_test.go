package datacatalog

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauv2"
	"github.com/jarcoal/httpmock"
	"github.com/reearth/reearthx/rerror"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFetcher(t *testing.T) {
	const base, prj = "", ""

	if base == "" || prj == "" {
		t.SkipNow()
	}
	f := lo.Must(NewFetcher(base))
	cmsres := lo.Must(f.Do(context.Background(), prj, FetcherDoOptions{}))
	res := cmsres.All()
	// item, _ := lo.Find(cmsres.Plateau, func(i PlateauItem) bool { return i.CityName == "" })
	// res := item.AllDataCatalogItems(item.IntermediateItem())
	// t.Log(string(lo.Must(json.MarshalIndent(res, "", "  "))))
	lo.Must0(os.WriteFile("datacatalog.json", lo.Must(json.MarshalIndent(res, "", "  ")), 0644))
}

func TestFetcher_Do(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/plateau", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "x"}},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/usecase", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "y"}},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/dataset", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "z"}},
		"totalCount": 1,
	})))

	ctx := context.Background()
	r, err := lo.Must(NewFetcher("https://example.com")).Do(ctx, "ppp", FetcherDoOptions{})
	assert.Equal(t, ResponseAll{
		Plateau: []plateauv2.CMSItem{{ID: "x"}},
		Usecase: []UsecaseItem{{ID: "y", Type: "ユースケース"}, {ID: "z"}},
	}, r)
	assert.NoError(t, err)
}

func TestFetcher_Do2(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/plateau", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusNotFound, map[string]any{
		"error": "not found",
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/usecase", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusNotFound, map[string]any{
		"error": "not found",
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/dataset", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "z"}},
		"totalCount": 1,
	})))

	ctx := context.Background()
	r, err := lo.Must(NewFetcher("https://example.com")).Do(ctx, "ppp", FetcherDoOptions{})
	assert.Equal(t, ResponseAll{
		Plateau: nil,
		Usecase: []UsecaseItem{{ID: "z"}},
	}, r)
	assert.NoError(t, err)
}

func TestFetcher_Do3(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/plateau", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusNotFound, map[string]any{
		"error": "not found",
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/usecase", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusNotFound, map[string]any{
		"error": "not found",
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/dataset", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusNotFound, map[string]any{
		"error": "not found",
	})))

	ctx := context.Background()
	r, err := lo.Must(NewFetcher("https://example.com")).Do(ctx, "ppp", FetcherDoOptions{})
	assert.Equal(t, rerror.ErrNotFound, err)
	assert.Empty(t, r)
}

func TestFetcher_Do4(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/plateau", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "x"}},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/usecase", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "y"}},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/ppp/dataset", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "z"}},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/subprj/plateau", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results": []any{
			map[string]any{"id": "a", "city_name": "xxx市"},
			map[string]any{"id": "b", "city_name": "yyy市"},
		},
		"totalCount": 1,
	})))
	httpmock.RegisterResponderWithQuery("GET", "https://example.com/api/p/subprj/dataset", "page=1&per_page=100", lo.Must(httpmock.NewJsonResponder(http.StatusOK, map[string]any{
		"results":    []any{map[string]any{"id": "c", "city_name": "xxx市"}},
		"totalCount": 1,
	})))

	ctx := context.Background()
	r, err := lo.Must(NewFetcher("https://example.com")).Do(ctx, "ppp", FetcherDoOptions{
		Subproject: "subprj",
		CityName:   "xxx市",
	})
	assert.Equal(t, ResponseAll{
		Plateau: []plateauv2.CMSItem{{ID: "x"}, {ID: "a", CityName: "xxx市"}},
		Usecase: []UsecaseItem{{ID: "y", Type: "ユースケース"}, {ID: "z"}, {ID: "c", CityName: "xxx市"}},
	}, r)
	assert.NoError(t, err)
}
