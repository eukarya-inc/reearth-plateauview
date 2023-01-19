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

	expected := "[{\"id\":\"df1\",\"type\":\"text\",\"value\":\"t\"}]"
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

func TestHandler_fetchTemplate(t *testing.T) {
	h := newHandler()
	modelID := "key2"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "[\"[{'hoge':'hoge'}]\"]"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, &cms.Items{
			Items: []cms.Item{
				{
					ID:     "a",
					Fields: []cms.Field{{ID: h.TemplateModelTemplateFieldID, Type: "text", Value: "[{'hoge':'hoge'}]"}},
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
	p := path.Join("/viz/aaa/templates/")
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
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

func TestHandler_createTemplateHandler(t *testing.T) {
	h := newHandler()
	modelID := "key2"
	httpmock.Activate()
	defer httpmock.Deactivate()

	expected := "{\"id\":\"tf1\",\"type\":\"TextArea\",\"value\":\"[{'hoge':'hoge'}]\"}"
	responder := func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, cms.Item{
			ID: modelID,
			Fields: []cms.Field{
				{ID: h.TemplateModelTemplateFieldID, Type: "TextArea", Value: "[{'hoge':'hoge'}]"},
				{ID: h.TemplateModelIDFieldID, Type: "Text", Value: "[{'hoge':'hoge'}]"},
			},
		},
		)
	}
	httpmock.RegisterResponder("POST", lo.Must(url.JoinPath(cmsHost, "/api/models/", modelID, "items")), responder)
	e := echo.New()
	p := path.Join("/viz/aaa/templates/")
	req := httptest.NewRequest(http.MethodGet, p, nil)
	req.Header.Set("Content-Type", "application/json")
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
