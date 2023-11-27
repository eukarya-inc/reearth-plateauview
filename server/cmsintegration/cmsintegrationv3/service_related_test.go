package cmsintegrationv3

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestConvertRelatedDataset(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://example.com/hoge_border.geojson",
		httpmock.NewStringResponder(200, `{}`))

	var updatedFields [][]*cms.Field
	var updatedMetadataFields [][]*cms.Field
	var uploaded []string
	var comments []string
	ctx := context.Background()
	c := &cmsMock{
		asset: func(ctx context.Context, id string) (*cms.Asset, error) {
			return &cms.Asset{
				URL: "https://example.com/hoge_border.geojson",
			}, nil
		},
		updateItem: func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			updatedFields = append(updatedFields, fields)
			updatedMetadataFields = append(updatedMetadataFields, metadataFields)
			return nil, nil
		},
		uploadAssetDirectly: func(ctx context.Context, prjectID, name string, r io.Reader) (string, error) {
			uploaded = append(uploaded, name)
			return "asset", nil
		},
		commentToItem: func(ctx context.Context, id string, comment string) error {
			comments = append(comments, comment)
			return nil
		},
	}
	s := &Services{CMS: c, HTTP: http.DefaultClient}
	item := &RelatedItem{
		Assets: map[string][]string{
			"border": {"border"},
		},
	}
	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Model: &cms.Model{
				Key: "plateau-related",
			},
			Item: item.CMSItem(),
		},
	}

	t.Run("sucess", func(t *testing.T) {
		updatedFields = nil
		updatedMetadataFields = nil
		uploaded = nil
		comments = nil

		err := convertRelatedDataset(ctx, s, w, item)
		assert.NoError(t, err)
		assert.Equal(t, [][]*cms.Field{
			nil,
			{
				{
					Key:   "border_conv",
					Type:  "asset",
					Value: []string{"asset"},
				},
			},
		}, updatedFields)
		assert.Equal(t, [][]*cms.Field{
			{
				{Key: "conv_status", Type: "select", Value: ConvertionStatusRunning},
			},
			{
				{Key: "conv_status", Type: "select", Value: ConvertionStatusSuccess},
			},
		}, updatedMetadataFields)
		assert.Equal(t, []string{"hoge_border.czml"}, uploaded)
		assert.Equal(t, []string{"変換を開始しました。", "変換に成功しました。"}, comments)
	})
}

func TestPackRelatedDataset(t *testing.T) {
	mockGeoJSON := func(name string) map[string]any {
		return map[string]any{
			"type": "FeatureCollection",
			"features": []any{
				map[string]any{
					"type": "Feature",
					"properties": map[string]any{
						"name": name,
					},
					"geometry": map[string]any{
						"type":        "Point",
						"coordinates": []any{0.0, 0.0},
					},
				},
			},
		}
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		`=~^https://example\.com/(.+)\.geojson`,
		func(req *http.Request) (*http.Response, error) {
			name := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewJsonResponse(200, mockGeoJSON(name))
		},
	)

	var updatedFields [][]*cms.Field
	var updatedMetadataFields [][]*cms.Field
	var uploadedData [][]byte
	var uploaded []string
	var comments []string
	ctx := context.Background()
	c := &cmsMock{
		getItem: func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			return (&CityItem{
				CityNameEn: "city",
				CityCode:   "code",
			}).CMSItem(), nil
		},
		asset: func(ctx context.Context, id string) (*cms.Asset, error) {
			return &cms.Asset{
				URL: fmt.Sprintf("https://example.com/%s.geojson", id),
			}, nil
		},
		updateItem: func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			updatedFields = append(updatedFields, fields)
			updatedMetadataFields = append(updatedMetadataFields, metadataFields)
			return nil, nil
		},
		uploadAssetDirectly: func(ctx context.Context, prjectID, name string, r io.Reader) (string, error) {
			uploaded = append(uploaded, name)
			b := bytes.NewBuffer(nil)
			_, _ = io.Copy(b, r)
			uploadedData = append(uploadedData, b.Bytes())
			return "asset", nil
		},
		commentToItem: func(ctx context.Context, id string, comment string) error {
			comments = append(comments, comment)
			return nil
		},
	}
	s := &Services{CMS: c, HTTP: http.DefaultClient}
	item := &RelatedItem{
		City: "city",
		Assets: map[string][]string{
			"shelter":         {"shelter"},
			"landmark":        {"landmark1", "landmark2"},
			"station":         {"station"},
			"park":            {"park"},
			"railway":         {"railway"},
			"emergency_route": {"emergency_route"},
			"border":          {"border"},
		},
	}
	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Model: &cms.Model{
				Key: "plateau-related",
			},
			Item: item.CMSItem(),
		},
	}

	t.Run("sucess", func(t *testing.T) {
		updatedFields = nil
		updatedMetadataFields = nil
		uploaded = nil
		comments = nil

		err := packRelatedDataset(ctx, s, w, item)
		assert.NoError(t, err)
		assert.Equal(t, [][]*cms.Field{
			nil,
			{
				{
					Key:   "merged",
					Type:  "asset",
					Value: "asset",
				},
			},
		}, updatedFields)
		assert.Equal(t, [][]*cms.Field{
			{
				{Key: "merge_status", Type: "select", Value: ConvertionStatusRunning},
			},
			{
				{Key: "merge_status", Type: "select", Value: ConvertionStatusSuccess},
			},
		}, updatedMetadataFields)
		assert.Equal(t, []string{"code_city_related.zip"}, uploaded)
		assert.Equal(t, []string{
			"G空間情報センター公開用zipファイルの作成を開始しました。",
			"G空間情報センター公開用zipファイルの作成が完了しました。",
		}, comments)

		zr, _ := zip.NewReader(bytes.NewReader(uploadedData[0]), int64(len(uploadedData[0])))
		assert.Equal(t, []string{
			"shelter.geojson",
			"park.geojson",
			"landmark1.geojson",
			"landmark2.geojson",
			"landmark.geojson",
			"station.geojson",
			"railway.geojson",
			"emergency_route.geojson",
			"border.geojson",
		}, lo.Map(zr.File, func(f *zip.File, _ int) string {
			return f.Name
		}))

		// assert landmark1.geojson
		zf := lo.Must(zr.Open("landmark1.geojson"))
		var ge map[string]any
		_ = json.NewDecoder(zf).Decode(&ge)
		assert.Equal(t, mockGeoJSON("landmark1"), ge)

		// assert landmark.geojson
		zf = lo.Must(zr.Open("landmark.geojson"))
		ge = nil
		_ = json.NewDecoder(zf).Decode(&ge)
		assert.Equal(t, map[string]any{
			"type": "FeatureCollection",
			"name": "code_city_landmark",
			"features": []any{
				map[string]any{
					"type":       "Feature",
					"properties": map[string]any{"name": "landmark1"},
					"geometry":   map[string]any{"type": "Point", "coordinates": []any{0.0, 0.0}},
				},
				map[string]any{
					"type":       "Feature",
					"properties": map[string]any{"name": "landmark2"},
					"geometry":   map[string]any{"type": "Point", "coordinates": []any{0.0, 0.0}},
				},
			},
		}, ge)
	})
}
