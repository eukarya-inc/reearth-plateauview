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
)

const key1 = "N03_001"
const key2 = "N03_004"
const dirPath = "govpolygondata"

type Handler struct {
	// e.g. "http://[::]:8080"
	gqlEndpoint string
	processor   *Processor
	httpClient  *http.Client
	lock        sync.RWMutex
	geojson     []byte
}

func New(gqlEndpoint string) *Handler {
	return &Handler{
		gqlEndpoint: gqlEndpoint,
		processor:   NewProcessor(dirPath, key1, key2),
		httpClient:  &http.Client{},
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

func (h *Handler) getCityNames(ctx context.Context) ([]string, error) {

	query := `
		{
			areas(input:{
				areaTypes: [CITY]
			}) {
				name
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
		return nil, err
	}

	req, err := http.NewRequest("POST", h.gqlEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData struct {
		Data struct {
			Areas []struct {
				Name       string `json:"name"`
				Prefecture struct {
					Name string `json:"name"`
				} `json:"prefecture"`
			} `json:"areas"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}

	cityNames := make([]string, len(responseData.Data.Areas))
	for i, city := range responseData.Data.Areas {
		cityNames[i] = city.Prefecture.Name + city.Name
	}

	return cityNames, nil
}
