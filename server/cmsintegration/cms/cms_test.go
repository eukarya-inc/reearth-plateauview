package cms

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var _ Interface = (*CMS)(nil)

func TestCMS(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	ctx := context.Background()

	// valid
	call := mockCMS("http://fme.example.com", "TOKEN")
	f := lo.Must(New("http://fme.example.com", "TOKEN"))
	assetID, err := f.UploadAsset(ctx, "aaa")
	assert.NoError(t, err)
	assert.Equal(t, "idid", assetID)
	assert.NoError(t, f.UpdateItem(ctx, "a", map[string]any{}))
	assert.NoError(t, f.Comment(ctx, "c", "comment"))
	assert.Equal(t, 1, call("POST /api/assets"))
	assert.Equal(t, 1, call("PATCH /api/items/a"))
	assert.Equal(t, 1, call("POST /api/threads/c/comments"))

	// invalid token
	httpmock.Reset()
	call = mockCMS("http://fme.example.com", "TOKEN")
	f = lo.Must(New("http://fme.example.com", "TOKEN2"))
	assetID, err = f.UploadAsset(ctx, "aaa")
	assert.ErrorContains(t, err, "failed to request: code=401")
	assert.Equal(t, "", assetID)
	assert.ErrorContains(t, f.UpdateItem(ctx, "a", map[string]any{}), "failed to request: code=401")
	assert.ErrorContains(t, f.Comment(ctx, "c", "comment"), "failed to request: code=401")
	assert.Equal(t, 1, call("POST /api/assets"))
	assert.Equal(t, 1, call("PATCH /api/items/a"))
	assert.Equal(t, 1, call("POST /api/threads/c/comments"))
}

func mockCMS(host, token string) func(string) int {
	responder := func(req *http.Request) (*http.Response, error) {
		if t := parseToken(req); t != token {
			return httpmock.NewJsonResponse(http.StatusUnauthorized, map[string]any{})
		}

		res := map[string]string{}
		p := req.URL.Path
		if p == "/api/assets" {
			res["id"] = "idid"
		}

		return httpmock.NewJsonResponse(http.StatusOK, res)
	}

	httpmock.RegisterResponder("PATCH", host+"/api/items/a", responder)
	httpmock.RegisterResponder("POST", host+"/api/assets", responder)
	httpmock.RegisterResponder("POST", host+"/api/threads/c/comments", responder)

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
