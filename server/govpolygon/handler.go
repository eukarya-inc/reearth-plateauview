package govpolygon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
	"github.com/rubenv/topojson"
)

const key1 = "N03_001"
const key2 = "N03_004"
const dirPath = "govpolygondata"

type Handler struct {
	// e.g. "http://[::]:8080"
	gqlEndpoint       string
	processor         *Processor
	httpClient        *http.Client
	lock              sync.RWMutex
	geojson           []byte
	topojson          []byte
	updateIfNotExists bool
}

func New(gqlEndpoint string, updateIfNotExists bool) *Handler {
	return &Handler{
		gqlEndpoint:       gqlEndpoint,
		processor:         NewProcessor(dirPath, key1, key2),
		httpClient:        http.DefaultClient,
		updateIfNotExists: updateIfNotExists,
	}
}

func (h *Handler) Route(g *echo.Group) *Handler {
	g.GET("/geojson", h.GetGeoJSON)
	g.GET("/topojson", h.GetTeopoJSON)
	// g.GET("/update", h.Update, errorLogger)
	return h
}

func (h *Handler) GetGeoJSON(c echo.Context) error {
	if h.updateIfNotExists && h.geojson == nil {
		if err := h.Update(c); err != nil {
			log.Errorfc(c.Request().Context(), "govpolygon: fail to init: %v", err)
		}
	}

	h.lock.RLock()
	defer h.lock.RUnlock()
	if h.geojson == nil {
		return c.JSON(http.StatusNotFound, "not found")
	}
	return c.JSONBlob(http.StatusOK, h.geojson)
}

func (h *Handler) GetTeopoJSON(c echo.Context) error {
	if h.updateIfNotExists && h.topojson == nil {
		if err := h.Update(c); err != nil {
			log.Errorfc(c.Request().Context(), "govpolygon: fail to init: %v", err)
		}
	}

	h.lock.RLock()
	defer h.lock.RUnlock()
	if h.topojson == nil {
		return c.JSON(http.StatusNotFound, "not found")
	}
	return c.JSONBlob(http.StatusOK, h.topojson)
}

func (h *Handler) Update(c echo.Context) error {
	initial := h.geojson == nil
	if initial {
		h.lock.Lock()
		defer h.lock.Unlock()
		if h.geojson != nil {
			return nil
		}
	}

	ctx := c.Request().Context()
	values, citycodem, err := h.getCityNames(ctx)
	if err != nil {
		return err
	}

	g, err := h.processor.ComputeGeoJSON(ctx, values, citycodem)
	if err != nil {
		return err
	}

	t := topojson.NewTopology(g, &topojson.TopologyOptions{})
	geojsonj, err := json.Marshal(g)
	if err != nil {
		return fmt.Errorf("failed to marshal geojson: %w", err)
	}
	topojsonj, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal topojson: %w", err)
	}

	if !initial {
		h.lock.Lock()
		defer h.lock.Unlock()
	}

	h.geojson = geojsonj
	h.topojson = topojsonj

	return nil
}

func (h *Handler) getCityNames(ctx context.Context) ([]string, map[string]string, error) {
	query := `
		{
			areas(input:{
				areaTypes: [CITY]
			}) {
				name
				code
				... on City {
					prefecture {
						name
					}
				}
			}
		}
	`

	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.gqlEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData struct {
		Data struct {
			Areas []struct {
				Name       string `json:"name"`
				Code       string `json:"code"`
				Prefecture struct {
					Name string `json:"name"`
				} `json:"prefecture"`
			} `json:"areas"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	m := map[string]string{}
	cityNames := make([]string, len(responseData.Data.Areas))
	for i, city := range responseData.Data.Areas {
		cityNames[i] = city.Prefecture.Name + city.Name
		m[cityNames[i]] = city.Code
	}

	return cityNames, m, nil
}

// func errorLogger(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		err := next(c)
// 		if err != nil {
// 			log.Errorfc(c.Request().Context(), "govpolygon: %v", err)
// 			if !c.Response().Committed {
// 				return c.JSON(http.StatusInternalServerError, map[string]any{"error": "internal"})
// 			}
// 		}
// 		return nil
// 	}
// }
