package visualizer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

// ctx := context.Background()

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
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("GET", lo.Must(url.JoinPath(cmsHost, "/api/items/", itemID)), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/data/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/data/:iid")
	ctx.SetParamNames("pid", "iid")
	ctx.SetParamValues("aaa", itemID)
	handler := h.getDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_getAllData(t *testing.T) {
	h := newHandler()
	modelID := "key1"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "{'items': [{'id':'a', 'fields': [{'id': 'f', 'Type': 'text', 'Value': 't'}],}],'page': 1,'perPage':50,'totalCount': 1}"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, &cms.Items{
			Items: []cms.Item{
				{
					ID:     "a",
					Fields: []cms.Field{{ID: h.DataModelDataFieldID, Type: "text", Value: "t"}},
				},
			},
			Page:       1,
			PerPage:    50,
			TotalCount: 1,
		},
		)
	}
	httpmock.RegisterResponder("GET", lo.Must(url.JoinPath(cmsHost, "/api/models/", modelID, "items")), responder)

	e := echo.New()
	p := path.Join("/viz/aaa/data/")
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	// TODO: recのbodyが空になってしまう
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/data/")
	ctx.SetParamNames("pid")
	ctx.SetParamValues("aaa")
	handler := h.getAllDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_createDataHandler(t *testing.T) {
	h := newHandler()
	modelID := "key1"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "[{'hoge':'hoge'}]"
	responder := func(req *http.Request) (*http.Response, error) {
		// token := getToken(req)
		// if cmsToken != token {
		// 	return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
		// }

		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: modelID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("POST", lo.Must(url.JoinPath(cmsHost, "/api/models/", modelID, "items")), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/data/")
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/data/")
	ctx.SetParamNames("pid")
	ctx.SetParamValues("aaa")
	handler := h.createDataHandler()
	err := handler(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_updateDataHandler(t *testing.T) {
	h := newHandler()
	itemID := "aaa"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "{'hoge':'hoge'}"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("PATCH", lo.Must(url.JoinPath(cmsHost, "/api/items/", itemID)), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/data/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/data/:iid")
	ctx.SetParamNames("pid", "iid")
	ctx.SetParamValues("aaa", itemID)
	handler := h.updateDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_deleteDataHandler(t *testing.T) {
	h := newHandler()
	itemID := "aaa"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "null"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("DELETE", lo.Must(url.JoinPath(cmsHost, "/api/items/", itemID)), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/data/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/data/:iid")
	ctx.SetParamNames("pid", "iid")
	ctx.SetParamValues("aaa", itemID)
	handler := h.deleteDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

/*
	func TestHandler_fetchTemplate(t *testing.T) {
		h := newHandler()
		modelID := "key2"
		httpmock.Activate()
		defer httpmock.Deactivate()

		expected := "[{'hoge':'hoge'}]"
		responder := func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
				ID: modelID,
				Fields: []cms.Field{
					{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
					{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
				},
			},
			)
		}
		httpmock.RegisterResponder("GET", lo.Must(url.JoinPath(cmsHost, "/api/models/", modelID, "items")), responder)
		e := echo.New()
		p := path.Join("/viz/aaa/templates/")
		req := httptest.NewRequest(http.MethodGet, p, nil)
		req.Header.Set("Content-Type", "application/json")
		// TODO: recのbodyが空になってしまう
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/viz/:pid/templates/")
		ctx.SetParamNames("pid")
		ctx.SetParamValues("aaa")
		handler := h.fetchTemplate()
		res := handler(ctx)
		assert.NoError(t, res)
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
	}
*/
func TestHandler_createTemplateHandler(t *testing.T) {
	h := newHandler()
	modelID := "key2"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "[{'hoge':'hoge'}]"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: modelID,
			Fields: []cms.Field{
				{ID: h.TemplateModelTemplateFieldID, Type: "TextArea", Value: expected},
				{ID: h.TemplateModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("POST", lo.Must(url.JoinPath(cmsHost, "/api/models/", modelID, "items")), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/templates/")
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	// TODO: recのbodyが空になってしまう
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/templates/")
	ctx.SetParamNames("pid")
	ctx.SetParamValues("aaa")
	handler := h.createTemplateHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_updateTemplateHandler(t *testing.T) {
	h := newHandler()
	itemID := "aaa"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "{'hoge':'hoge'}"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("PATCH", lo.Must(url.JoinPath(cmsHost, "/api/items/", itemID)), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/templates/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/templates/:iid")
	ctx.SetParamNames("pid", "iid")
	ctx.SetParamValues("aaa", itemID)
	handler := h.updateDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

func TestHandler_deleteTemplateHandler(t *testing.T) {
	h := newHandler()
	itemID := "aaa"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "null"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: itemID,
			Fields: []cms.Field{
				{ID: h.DataModelDataFieldID, Type: "TextArea", Value: expected},
				{ID: h.DataModelIDFieldID, Type: "Text", Value: expected},
			},
		},
		)
	}
	httpmock.RegisterResponder("DELETE", lo.Must(url.JoinPath(cmsHost, "/api/items/", itemID)), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/templates/", itemID)
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/viz/:pid/templates/:iid")
	ctx.SetParamNames("pid", "iid")
	ctx.SetParamValues("aaa", itemID)
	handler := h.deleteDataHandler()
	res := handler(ctx)
	assert.NoError(t, res)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, expected, strings.Trim(strings.TrimSpace(rec.Body.String()), "\""))
}

// func mockCMS(host, token string) func(string) int {
// 	responder := func(req *http.Request) (*http.Response, error) {
// 		if t := parseToken(req); t != token {
// 			return httpmock.NewJsonResponse(http.StatusUnauthorized, "unauthorized")
// 		}

// 		if req.Header.Get("Content-Type") != "application/json" {
// 			return httpmock.NewJsonResponse(http.StatusUnsupportedMediaType, "unsupported media type")
// 		}

// 		res := map[string]any{}
// 		p := req.URL.Path
// 		if req.Method == "POST" && p == "/api/projects/ppp/assets" {
// 			res["id"] = "idid"
// 		} else if req.Method == "GET" && p == "/api/assets/a" {
// 			res["id"] = "a"
// 			res["url"] = "url"
// 		} else if req.Method == "POST" && p == "/api/models/a/items" || p == "/api/items/a" {
// 			res["id"] = "a"
// 			// TODO: fields あとで変える
// 			res["fields"] = []map[string]string{{"id": "f", "type": "text", "value": "t"}}
// 		} else if req.Method == "PATCH" && p == "/api/models/a/items/" {
// 			//TDOO: PATCHのときの処理を書く
// 		} else if req.Method == "DELETE" && p == "/api/models/a/items" {
// 			res = nil
// 		}

// 		return httpmock.NewJsonResponse(http.StatusOK, res)
// 	}

// 	httpmock.RegisterResponder("GET", host+"/api/items/a", responder)
// 	httpmock.RegisterResponder("PATCH", host+"/api/items/a", responder)
// 	httpmock.RegisterResponder("POST", host+"/api/models/a/items", responder)
// 	httpmock.RegisterResponder("POST", host+"/api/projects/ppp/assets", responder)
// 	httpmock.RegisterResponder("POST", host+"/api/assets/c/comments", responder)
// 	httpmock.RegisterResponder("GET", host+"/api/assets/a", responder)

// 	return func(p string) int {
// 		b, a, _ := strings.Cut(p, " ")
// 		return httpmock.GetCallCountInfo()[b+" "+host+a]
// 	}
// }

func getToken(r *http.Request) string {
	aut := r.Header.Get("Authorization")
	_, token, found := strings.Cut(aut, "Bearer ")
	if !found {
		return ""
	}
	return token
}
