package cms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/reearth/reearthx/log"
)

type Interface interface {
	GetItem(ctx context.Context, itemID string) (*Item, error)
<<<<<<< HEAD
	GetItems(ctx context.Context, modelID string) ([]*Item, error)
=======
	GetItems(ctx context.Context, modelID string) (*Items, error)
	GetItemsByKey(ctx context.Context, projectIDOrAlias, modelIDOrKey string) (*Items, error)
>>>>>>> main
	CreateItem(ctx context.Context, modelID string, fields []Field) (*Item, error)
	UpdateItem(ctx context.Context, itemID string, fields []Field) (*Item, error)
	DeleteItem(ctx context.Context, itemID string) error
	Asset(ctx context.Context, id string) (*Asset, error)
	UploadAsset(ctx context.Context, projectID, url string) (string, error)
	UploadAssetDirectly(ctx context.Context, projectID, name string, data io.Reader) (string, error)
	CommentToItem(ctx context.Context, assetID, content string) error
	CommentToAsset(ctx context.Context, assetID, content string) error
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

func (c *CMS) GetItems(ctx context.Context, modelID string) ([]*Item, error) {
	b, err := c.send(ctx, http.MethodGet, []string{"api", "models", modelID, "items"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send an request: %w", err)
	}
	byte, err := io.ReadAll(b)
	if err != nil {
		return nil, fmt.Errorf("occur an unexpected EOF error: %w", err)
	}

	defer func() { _ = b.Close() }()

	var items Items
	if err := json.Unmarshal(byte, &items); err != nil {
		return nil, err
	}
	return items.Items, nil
}

func (c *CMS) GetItem(ctx context.Context, itemID string) (*Item, error) {
	b, err := c.send(ctx, http.MethodGet, []string{"api", "items", itemID}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get an item: %w", err)
	}
	defer func() { _ = b.Close() }()

	item := &Item{}
	if err := json.NewDecoder(b).Decode(item); err != nil {
		return nil, fmt.Errorf("failed to parse an item: %w", err)
	}

	return item, nil
}

func (c *CMS) GetItems(ctx context.Context, modelID string) (*Items, error) {
	b, err := c.send(ctx, http.MethodGet, []string{"api", "models", modelID, "items"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer func() { _ = b.Close() }()

	items := &Items{}
	if err := json.NewDecoder(b).Decode(items); err != nil {
		return nil, fmt.Errorf("failed to parse items: %w", err)
	}

	return items, nil
}

func (c *CMS) GetItemsByKey(ctx context.Context, projectIDOrAlias, modelIDOrAlias string) (*Items, error) {
	b, err := c.send(ctx, http.MethodGet, []string{"api", "projects", projectIDOrAlias, "models", modelIDOrAlias, "items"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer func() { _ = b.Close() }()

	items := &Items{}
	if err := json.NewDecoder(b).Decode(items); err != nil {
		return nil, fmt.Errorf("failed to parse items: %w", err)
	}

	return items, nil
}

func (c *CMS) CreateItem(ctx context.Context, modelID string, fields []Field) (*Item, error) {
	rb := map[string]any{
		"fields": fields,
	}

	b, err := c.send(ctx, http.MethodPost, []string{"api", "models", modelID, "items"}, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to create an item: %w", err)
	}
	defer func() { _ = b.Close() }()

	item := &Item{}
	if err := json.NewDecoder(b).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to parse an item: %w", err)
	}

	return item, nil
}

func (c *CMS) UpdateItem(ctx context.Context, itemID string, fields []Field) (*Item, error) {
	rb := map[string]any{
		"fields": fields,
	}

	b, err := c.send(ctx, http.MethodPatch, []string{"api", "items", itemID}, rb)
	if err != nil {
		return nil, fmt.Errorf("failed to update an item: %w", err)
	}
	defer func() { _ = b.Close() }()

	item := &Item{}
	if err := json.NewDecoder(b).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to parse an item: %w", err)
	}

	return item, nil
}

func (c *CMS) DeleteItem(ctx context.Context, itemID string) error {
	b, err := c.send(ctx, http.MethodDelete, []string{"api", "items", itemID}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete an item: %w", err)
	}

	defer func() { _ = b.Close() }()
	return nil
}

func (c *CMS) UploadAsset(ctx context.Context, projectID, url string) (string, error) {
	rb := map[string]string{
		"url": url,
	}

	b, err2 := c.send(ctx, http.MethodPost, []string{"api", "projects", projectID, "assets"}, rb)
	if err2 != nil {
		log.Errorf("cms: upload asset: failed to upload an asset: %s", err2)
		return "", fmt.Errorf("failed to upload an asset: %w", err2)
	}

	defer func() { _ = b.Close() }()

	body, err2 := io.ReadAll(b)
	if err2 != nil {
		return "", fmt.Errorf("failed to read body: %w", err2)
	}

	type res struct {
		ID string `json:"id"`
	}

	r := &res{}
	if err2 := json.Unmarshal(body, &r); err2 != nil {
		return "", fmt.Errorf("failed to parse body: %w", err2)
	}

	return r.ID, nil
}

func (c *CMS) UploadAssetDirectly(ctx context.Context, projectID, name string, data io.Reader) (string, error) {
	// TODO
	return "", errors.New("not implemented yet")
}

func (c *CMS) Asset(ctx context.Context, assetID string) (*Asset, error) {
	b, err := c.send(ctx, http.MethodGet, []string{"api", "assets", assetID}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get an asset: %w", err)
	}
	defer func() { _ = b.Close() }()

	a := &Asset{}
	if err := json.NewDecoder(b).Decode(a); err != nil {
		return nil, fmt.Errorf("failed to parse an asset: %w", err)
	}

	return a, nil
}

func (c *CMS) CommentToItem(ctx context.Context, itemID, content string) error {
	rb := map[string]string{
		"content": content,
	}

	b, err := c.send(ctx, http.MethodPost, []string{"api", "items", itemID, "comments"}, rb)
	if err != nil {
		return fmt.Errorf("failed to comment to item %s: %w", itemID, err)
	}
	defer func() { _ = b.Close() }()

	return nil
}

func (c *CMS) CommentToAsset(ctx context.Context, assetID, content string) error {
	rb := map[string]string{
		"content": content,
	}

	b, err := c.send(ctx, http.MethodPost, []string{"api", "assets", assetID, "comments"}, rb)
	if err != nil {
		return fmt.Errorf("failed to comment to asset %s: %w", assetID, err)
	}
	defer func() { _ = b.Close() }()

	return nil
}

func (c *CMS) send(ctx context.Context, m string, p []string, body any) (io.ReadCloser, error) {
	req, err := c.request(ctx, m, p, body)
	if err != nil {
		return nil, err
	}

	log.Infof("CMS: request: %s %s body=%+v", req.Method, req.URL, body)

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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return req, nil
}
<<<<<<< HEAD

type Asset struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Item struct {
	ID     string  `json:"id"`
	Fields []Field `json:"fields"`
}

type Items struct {
	Items      []*Item `json:"items"`
	TotalCount int64   `json:"totalCount"`
}

func (i Item) Field(id string) *Field {
	f, ok := lo.Find(i.Fields, func(f Field) bool { return f.ID == id })
	if ok {
		return &f
	}
	return nil
}
=======
>>>>>>> main
