package visualizer

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const (
	cmsHost                      = "https://api.cms.test.reearth.dev"
	cmsToken                     = "token"
	dataModelKey                 = "key1"
	templateModelKey             = "key2"
	dataModelDataFieldID         = "df1"
	dataModelIDFieldID           = "df2"
	templateModelTemplateFieldID = "tf1"
	templateModelIDFieldID       = "tf2"
)

// func TestHandler(t *testing.T) {
// httpmock.Activate()
// defer httpmock.Deactivate()
// mockCMS(t)

//  e := echo.New()

//	// ctx := context.Background()

//	r := httptest.NewRequest("GET", "/viz/aaa", nil)
//	w := httptest.NewRecorder()
//	e.ServeHTTP(w, r)
//	assert.Equal(t, http.StatusOK, w.Code)
//}

func newHandler() *Handler {
	CMS := lo.Must(cms.New(cmsHost, cmsToken))
	return &Handler{
		CMS:                          CMS,
		DataModelKey:                 dataModelKey,
		DataModelDataFieldID:         dataModelDataFieldID,
		DataModelIDFieldID:           dataModelIDFieldID,
		TemplateModelKey:             templateModelKey,
		TemplateModelTemplateFieldID: templateModelTemplateFieldID,
		TemplateModelIDFieldID:       templateModelIDFieldID,
	}
}

func TestHandler_getData(t *testing.T) {
	h := newHandler()
	itemID := "aaa"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "{'hoge':'hoge'}"
	//Mockでやりたいこと: dataのITEMを返してほしい
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("GET", path.Join(cmsHost, "/api/items/", itemID), responder)
	//テストしたいこと: CMSからdataが返ってくる想定のもと、仕様どおりにデータを返せるかどうか？
	e := echo.New()
	p := path.Join(cmsHost, "/viz/aaa/data/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// e.ServeHTTP(rec, req)
	assert.NoError(t, h.getDataHandler()(c))
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, rec.Body.String())
}

func mockCMS(host, token string) func(string) int {
	responder := func(req *http.Request) (*http.Response, error) {
		if t := parseToken(req); t != token {
			return httpmock.NewJsonResponse(http.StatusUnauthorized, "unauthorized")
		}

		if req.Header.Get("Content-Type") != "application/json" {
			return httpmock.NewJsonResponse(http.StatusUnsupportedMediaType, "unsupported media type")
		}

		res := map[string]any{}
		p := req.URL.Path
		if req.Method == "POST" && p == "/api/projects/ppp/assets" {
			res["id"] = "idid"
		} else if req.Method == "GET" && p == "/api/assets/a" {
			res["id"] = "a"
			res["url"] = "url"
		} else if req.Method == "POST" && p == "/api/models/a/items" || p == "/api/items/a" {
			res["id"] = "a"
			// TODO: fields あとで変える
			res["fields"] = []map[string]string{{"id": "f", "type": "text", "value": "t"}}
		} else if req.Method == "PATCH" && p == "/api/models/a/items/" {
			//TDOO: PATCHのときの処理を書く
		} else if req.Method == "DELETE" && p == "/api/models/a/items" {
			res = nil
		}

		return httpmock.NewJsonResponse(http.StatusOK, res)
	}

	httpmock.RegisterResponder("GET", host+"/api/items/a", responder)
	httpmock.RegisterResponder("PATCH", host+"/api/items/a", responder)
	httpmock.RegisterResponder("POST", host+"/api/models/a/items", responder)
	httpmock.RegisterResponder("POST", host+"/api/projects/ppp/assets", responder)
	httpmock.RegisterResponder("POST", host+"/api/assets/c/comments", responder)
	httpmock.RegisterResponder("GET", host+"/api/assets/a", responder)

	return func(p string) int {
		b, a, _ := strings.Cut(p, " ")
		return httpmock.GetCallCountInfo()[b+" "+host+a]
	}
}

func parseToken(r *http.Request) string {
	aut := r.Header.Get("Authorization")
	_, token, found := strings.Cut(aut, "Bearer ")
	if !found {
		return ""
	}
	return token
}
