package govpolygon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	url := ""
	if url == "" {
		t.Skip("skipping test; no URL provided")
	}
	h := New(url)

	ctx := context.Background()
	assert.NoError(t, h.Update(ctx))

	e := echo.New()
	r := httptest.NewRequest(http.MethodGet, "", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	geojson := h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body.String())

	t.Log(geojson)
}

func TestProcessor(t *testing.T) {
	p := &Processor{
		dirpath: "govpolygondata",
		key1:    key1,
		key2:    key2,
	}

	assert.NoError(t, p.Init())
	ctx := context.Background()
	values := []string{"東京都千代田区"}
	geojson, err := p.ComputeGeoJSON(ctx, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, geojson)
}
