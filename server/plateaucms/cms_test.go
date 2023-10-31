package plateaucms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/rerror"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const (
	testCMSHost            = "https://example.com"
	testCMSToken           = "token"
	testCMSProject         = "prj"
	testSidebarAccessToken = "access_token"
	testModelKey           = "model1"
	testModelKey2          = "model2"
)

func TestHandler_Metadata(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()
	mockCMS(t)
	h := newHandler()

	toekn, err := h.Metadata(context.Background(), "prjprj")
	assert.NoError(t, err)
	assert.Equal(t, Metadata{
		ProjectAlias:       "prjprj",
		CMSAPIKey:          "token!",
		SidebarAccessToken: "ac",
	}, toekn)

	toekn, err = h.Metadata(context.Background(), "prjprj!")
	assert.Equal(t, rerror.ErrNotFound, err)
	assert.Empty(t, toekn)
}

func TestHandler_LastModified(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()
	lastModified := time.Date(2022, time.April, 1, 0, 0, 0, 0, time.Local)
	lastModified2 := time.Date(2022, time.April, 2, 0, 0, 0, 0, time.Local)

	httpmock.RegisterResponder(
		"GET",
		lo.Must(url.JoinPath(testCMSHost, "api", "projects", testCMSProject, "models", testModelKey)),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, &cms.Model{LastModified: lastModified}),
	)
	httpmock.RegisterResponder(
		"GET",
		lo.Must(url.JoinPath(testCMSHost, "api", "projects", testCMSProject, "models", testModelKey2)),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, &cms.Model{LastModified: lastModified2}),
	)
	h := newHandler()
	cms := lo.Must(cms.New(testCMSHost, testCMSToken))

	e := echo.New()

	// no If-Modified-Since
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), cmsContextKey{}, cms))
	w := httptest.NewRecorder()
	hit, err := h.LastModified(e.NewContext(r, w), testCMSProject, testModelKey, testModelKey2)
	assert.NoError(t, err)
	assert.False(t, hit)
	assert.Equal(t, lastModified2.Format(time.RFC1123), w.Header().Get(echo.HeaderLastModified))

	// If-Modified-Since
	r = httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), cmsContextKey{}, cms))
	r.Header.Set(echo.HeaderIfModifiedSince, lastModified2.Format(time.RFC1123))
	w = httptest.NewRecorder()
	hit, err = newHandler().LastModified(e.NewContext(r, w), testCMSProject, testModelKey, testModelKey2)
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, http.StatusNotModified, w.Result().StatusCode)
	assert.Equal(t, lastModified2.Format(time.RFC1123), w.Header().Get(echo.HeaderLastModified))

	// expired If-Modified-Since
	r = httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), cmsContextKey{}, cms))
	r.Header.Set(echo.HeaderIfModifiedSince, lastModified.Format(time.RFC1123))
	w = httptest.NewRecorder()
	hit, err = newHandler().LastModified(e.NewContext(r, w), testCMSProject, testModelKey, testModelKey2)
	assert.NoError(t, err)
	assert.False(t, hit)
	assert.Equal(t, lastModified2.Format(time.RFC1123), w.Header().Get(echo.HeaderLastModified))
}

func newHandler() *CMS {
	return &CMS{
		cmsbase:         testCMSHost,
		cmsTokenProject: tokenProject,
		cmsMain:         lo.Must(cms.New(testCMSHost, testCMSToken)),
	}
}

func mockCMS(t *testing.T) {
	t.Helper()

	httpmock.RegisterResponder(
		"GET",
		lo.Must(url.JoinPath(testCMSHost, "api", "projects", tokenProject, "models", tokenModel, "items")),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, cms.Items{
			PerPage:    1,
			Page:       1,
			TotalCount: 1,
			Items: []cms.Item{
				{
					ID: "1",
					Fields: []*cms.Field{
						{Key: tokenProjectField, Value: testCMSProject},
						{Key: "cms_apikey", Value: testCMSToken},
						{Key: "sidebar_access_token", Value: testSidebarAccessToken},
					},
				},
				{
					ID: "2",
					Fields: []*cms.Field{
						{Key: tokenProjectField, Value: "prjprj"},
						{Key: "cms_apikey", Value: "token!"},
						{Key: "sidebar_access_token", Value: "ac"},
					},
				},
			},
		}),
	)
}
