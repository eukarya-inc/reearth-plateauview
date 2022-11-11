package cms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Interface interface {
	UploadAsset(ctx context.Context, r io.Reader) (string, error)
	UpdateItem(ctx context.Context, itemID string, fields map[string]any) error
	Comment(ctx context.Context, assetID, content string) error
}

type CMS struct {
	base   *url.URL
	token  string
	client *http.Client
}

func New(base, token string) (*CMS, error) {
	b, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}

	return &CMS{
		base:   b,
		token:  token,
		client: http.DefaultClient,
	}, nil
}

func (c *CMS) UploadAsset(ctx context.Context, r io.Reader) (string, error) {
	rb := map[string]string{}

	b, err := c.send(ctx, http.MethodPost, []string{"api", "assets"}, rb)
	if err != nil {
		return "", err
	}
	defer func() { _ = b.Close() }()

	var res map[string]any
	if err := json.NewDecoder(b).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to parse body: %w", err)
	}

	// do

	return "", nil
}

func (c *CMS) UpdateItem(ctx context.Context, itemID string, fields map[string]any) error {
	rb := map[string]any{
		"fields": fields,
	}

	b, err := c.send(ctx, http.MethodPatch, []string{"api", "items", itemID}, rb)
	if err != nil {
		return err
	}
	defer func() { _ = b.Close() }()

	var res map[string]any
	if err := json.NewDecoder(b).Decode(&res); err != nil {
		return fmt.Errorf("failed to parse body: %w", err)
	}

	// do

	return nil
}

func (c *CMS) Comment(ctx context.Context, assetID, content string) error {
	rb := map[string]string{
		"content": content,
	}

	b, err := c.send(ctx, http.MethodPost, []string{"api", "threads", assetID, "comments"}, rb)
	if err != nil {
		return err
	}
	defer func() { _ = b.Close() }()

	return nil
}

func (c *CMS) send(ctx context.Context, m string, p []string, body any) (io.ReadCloser, error) {
	req, err := c.request(ctx, m, p, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		defer func() {
			_ = res.Body.Close()
		}()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		return nil, fmt.Errorf("failed to request: code=%d, body=%s", res.StatusCode, b)
	}

	return res.Body, nil
}

func (c *CMS) request(ctx context.Context, m string, p []string, body any) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		bb, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		b = bytes.NewReader(bb)
	}

	req, err := http.NewRequestWithContext(ctx, m, c.base.JoinPath(p...).String(), b)
	if err != nil {
		return nil, fmt.Errorf("failed to init request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return req, nil
}
