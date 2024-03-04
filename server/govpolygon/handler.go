package govpolygon

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

const key1 = "N03_001"
const key2 = "N03_004"
const dirPath = "govpolygondata"

type Handler struct {
	// e.g. "http://[::]:8080"
	gqlEndpoint string
	processor   *Processor
	lock        sync.RWMutex
	geojson     []byte
}

func New(gqlEndpoint string) *Handler {
	return &Handler{
		gqlEndpoint: gqlEndpoint,
		processor:   NewProcessor(dirPath, key1, key2),
	}
}

func (h *Handler) Get(c echo.Context) error {
	if h.geojson == nil {
		return c.JSON(http.StatusNotFound, "not found")
	}

	h.lock.RLock()
	defer h.lock.RUnlock()
	return c.JSONBlob(http.StatusOK, h.geojson)
}

func (h *Handler) Update(ctx context.Context) error {
	values, err := h.getCityNames(ctx)
	if err != nil {
		return err
	}

	g, err := h.processor.ComputeGeoJSON(ctx, values)
	if err != nil {
		return err
	}

	h.lock.Lock()
	h.geojson = g
	h.lock.Unlock()

	return nil
}

/*
http://[::]:8080/datacatalog/graphql

{
  areas(input:{
    areaTypes: [CITY]
  }) {
    id
    name
    code
    ... on City {
      prefecture {
        name
      }
    }
  }
}
*/

func (h *Handler) getCityNames(ctx context.Context) ([]string, error) {
	// TODO e.g. ["東京都千代田区", "東京都世田谷区", ...]

	panic("not implemented")
}
