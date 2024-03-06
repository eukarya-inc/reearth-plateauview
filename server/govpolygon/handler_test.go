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

	e := echo.New()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	assert.NoError(t, h.Update(c))
	assert.NoError(t, h.GetGeoJSON(c))

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.NotEmpty(t, body)

	t.Log(body)
}

func TestProcessor(t *testing.T) {
	p := &Processor{
		dirpath: "govpolygondata",
		key1:    key1,
		key2:    key2,
	}

	ctx := context.Background()
	values := []string{"東京都千代田区"}
	geojson, err := p.ComputeGeoJSON(ctx, values, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, geojson)
}
