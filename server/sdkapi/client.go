package sdkapi

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type listResponse struct {
	Results Items `json:"results"`
}

type Client struct {
	c       *http.Client
	base    string
	project string
	model   string
}

func NewClient(c *http.Client, base, project, model string) *Client {
	if c == nil {
		c = http.DefaultClient
	}
	return &Client{
		c:       c,
		base:    base,
		project: project,
		model:   model,
	}
}

func (c *Client) GetItems(ctx context.Context) (Items, error) {
	u, err := url.JoinPath(c.base, "api", "p", c.project, c.model)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	var r listResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.Results, nil
}

func (c *Client) GetItem(ctx context.Context, id string) (*Item, error) {
	u, err := url.JoinPath(c.base, "api", "p", c.project, c.model, id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	var r Item
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *Client) GetMaxLOD(u string) (MaxLODColumns, error) {
	res, err := c.c.Get(u)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	r := csv.NewReader(res.Body)
	r.ReuseRecord = true
	var results MaxLODColumns
	for {
		c, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to read csv: %w", err)
		}

		if err == io.EOF {
			break
		}

		if len(c) != 3 || c[0] == "code" {
			continue
		}

		results = append(results, MaxLODColumn{
			Code:   c[0],
			Type:   c[1],
			MaxLOD: c[2],
		})
	}

	return results, nil
}
