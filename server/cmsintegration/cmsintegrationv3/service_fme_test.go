package cmsintegrationv3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSendRequestToFME(t *testing.T) {
	ctx := context.Background()
	c := &cmsMock{}
	f := &fmeMock{}
	s := &Services{
		CMS: c,
		FME: f,
	}
	conf := &Config{
		Secret: "secret",
	}
	baseItem := &cms.Item{
		ID:             "itemID",
		MetadataItemID: lo.ToPtr("metadataItemID"),
		Fields: []*cms.Field{
			{
				Key:   "citygml",
				Value: "citygmlID",
			},
			{
				Key:   "city",
				Value: "cityID",
			},
		},
	}
	cityItem := &cms.Item{
		ID: "cityID",
		Fields: []*cms.Field{
			{
				Key:   "codelists",
				Value: "codelistID",
			},
			{
				Key:   "citygml",
				Value: "citygmlID",
			},
		},
	}
	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Model: &cms.Model{
				Key: "modelKey",
			},
			Item: baseItem,
		},
	}

	// Test case 1: no metadataItemID and originalItemID
	_ = getLogs(t, func() {
		item := *baseItem
		item.MetadataItemID = nil
		item.OriginalItemID = nil
		w.ItemData.Item = &item

		c.reset()

		err := sendRequestToFME(ctx, s, conf, w)
		assert.ErrorContains(t, err, "invalid webhook payload")
	})

	// Test case 2: already converted
	log := getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
					Fields: []*cms.Field{
						{
							Key:   "conv_status",
							Value: string(ConvertionStatusSuccess),
						},
					},
				}, nil
			}
			return nil, fmt.Errorf("failed to get item")
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.NoError(t, err)
	})
	assert.Contains(t, log, "already converted")

	// Test case 3: skip convert
	log = getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
					Fields: []*cms.Field{
						{
							Key:   "skip_convert",
							Value: true,
						},
						{
							Key:   "skip_qc",
							Value: true,
						},
					},
				}, nil
			}
			return nil, fmt.Errorf("failed to get item")
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.NoError(t, err)
	})
	assert.Contains(t, log, "skip convert")

	// Test case 4: failed to get citygml asset
	_ = getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
				}, nil
			}
			return cityItem, nil
		}
		c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			return nil, nil
		}
		c.asset = func(ctx context.Context, id string) (*cms.Asset, error) {
			return nil, fmt.Errorf("failed to get citygml asset")
		}
		c.commentToItem = func(ctx context.Context, assetID, content string) error {
			assert.Contains(t, content, "CityGMLが見つかりません。")
			return nil
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.ErrorContains(t, err, "failed to get citygml asset")
	})

	// Test case 5: failed to get city item
	_ = getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
				}, nil
			}
			return nil, fmt.Errorf("failed to get city item")
		}
		c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			return nil, nil
		}
		c.asset = func(ctx context.Context, id string) (*cms.Asset, error) {
			return &cms.Asset{
				ID: "citygmlID",
			}, nil
		}
		c.commentToItem = func(ctx context.Context, assetID, content string) error {
			assert.Contains(t, content, "都市アイテムが見つかりません。")
			return nil
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.ErrorContains(t, err, "failed to get city item")
	})

	// Test case 6: failed to get codelist asset
	_ = getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
				}, nil
			}
			return cityItem, nil
		}
		c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			return nil, nil
		}
		c.asset = func(ctx context.Context, id string) (*cms.Asset, error) {
			if id == "citygmlID" {
				return &cms.Asset{
					ID: "citygmlID",
				}, nil
			}
			return nil, fmt.Errorf("failed to get codelist asset")
		}
		c.commentToItem = func(ctx context.Context, assetID, content string) error {
			assert.Contains(t, content, "コードリストが見つかりません。")
			return nil
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.ErrorContains(t, err, "failed to get codelist asset")
	})

	// Test case 7: success
	_ = getLogs(t, func() {
		item := *baseItem
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
				}, nil
			}
			return cityItem, nil
		}
		c.asset = func(ctx context.Context, id string) (*cms.Asset, error) {
			if id == "citygmlID" {
				return &cms.Asset{
					ID: "citygmlID",
				}, nil
			}
			return &cms.Asset{
				ID: "codelistID",
			}, nil
		}
		c.uploadAsset = func(ctx context.Context, projectID, url string) (string, error) {
			return "asset", nil
		}
		c.uploadAssetDirectly = func(ctx context.Context, projectID, name string, r io.Reader) (string, error) {
			return "assetd", nil
		}
		c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			assert.Equal(t, "conv_status", metadataFields[0].Key)
			assert.Equal(t, ConvertionStatusRunning, metadataFields[0].Value)
			return nil, nil
		}
		c.commentToItem = func(ctx context.Context, assetID, content string) error {
			assert.Contains(t, content, "品質検査・変換を開始しました。")
			return nil
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.NoError(t, err)
	})

	// Test case 8: success with metadata item
	_ = getLogs(t, func() {
		item := *baseItem
		item.ID = "metadataItemID"
		item.MetadataItemID = nil
		item.OriginalItemID = lo.ToPtr("itemID")
		w.ItemData.Item = &item

		c.reset()
		c.getItem = func(ctx context.Context, id string, asset bool) (*cms.Item, error) {
			if id == "itemID" {
				i := *baseItem
				return &i, nil
			}
			if id == "metadataItemID" {
				return &cms.Item{
					ID: "metadataItemID",
				}, nil
			}
			return cityItem, nil
		}
		c.asset = func(ctx context.Context, id string) (*cms.Asset, error) {
			if id == "citygmlID" {
				return &cms.Asset{
					ID: "citygmlID",
				}, nil
			}
			return &cms.Asset{
				ID: "codelistID",
			}, nil
		}
		c.uploadAsset = func(ctx context.Context, projectID, url string) (string, error) {
			return "asset", nil
		}
		c.uploadAssetDirectly = func(ctx context.Context, projectID, name string, r io.Reader) (string, error) {
			return "assetd", nil
		}
		c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
			assert.Equal(t, "conv_status", metadataFields[0].Key)
			assert.Equal(t, ConvertionStatusRunning, metadataFields[0].Value)
			return nil, nil
		}
		c.commentToItem = func(ctx context.Context, assetID, content string) error {
			assert.Contains(t, content, "品質検査・変換を開始しました。")
			return nil
		}

		err := sendRequestToFME(ctx, s, conf, w)
		assert.NoError(t, err)
	})
}

func TestReceiveResultFromFME(t *testing.T) {
	ctx := context.Background()
	c := &cmsMock{}
	s := &Services{
		CMS: c,
	}
	conf := &Config{
		Secret: "secret",
	}
	res := &fmeResult{
		Type: "conv",
		ID: fmeID{
			ItemID:      "itemID",
			ProjectID:   "projectID",
			FeatureType: "bldg",
			Type:        "qc_conv",
		}.String("secret"),
		LogURL: "log_ok",
		Results: map[string]any{
			"_dic":    "dic",
			"_maxlod": "maxlod",
			"bldg":    "bldg",
		},
		Status: "success",
	}

	// test case 1: success
	c.reset()
	uploaded := []string{}
	c.uploadAsset = func(ctx context.Context, projectID, url string) (string, error) {
		assert.Equal(t, projectID, "projectID")
		uploaded = append(uploaded, url)
		return url, nil
	}
	c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
		assert.Equal(t, id, "itemID")
		assert.Equal(t, []*cms.Field{
			{
				Key:   "data",
				Type:  "asset",
				Value: []string{"bldg"},
			},
			{
				Key:   "qc_result",
				Type:  "asset",
				Value: "log_ok",
			},
			{
				Key:   "dic",
				Type:  "asset",
				Value: "dic",
			},
			{
				Key:   "maxlod",
				Type:  "asset",
				Value: "maxlod",
			},
		}, fields)
		assert.Equal(t, []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "select",
				Value: ConvertionStatusSuccess,
			},
			{
				Key:   "qc_status",
				Type:  "select",
				Value: ConvertionStatusSuccess,
			},
		}, metadataFields)
		return nil, nil
	}
	c.commentToItem = func(ctx context.Context, assetID, content string) error {
		assert.Contains(t, content, "品質検査・変換が完了しました。")
		return nil
	}

	err := receiveResultFromFME(ctx, s, conf, *res)
	assert.NoError(t, err)
	assert.Equal(t, []string{"bldg", "dic", "maxlod", "log_ok"}, uploaded)

	// test case 2: invalid id
	r := *res
	r.ID = "invalid"
	err = receiveResultFromFME(ctx, s, conf, r)
	assert.ErrorContains(t, err, "invalid id")

	// test case 3: failed convert
	r = *res
	r.Status = "error"
	r.LogURL = "log"
	c.reset()
	c.updateItem = func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
		assert.Equal(t, []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "select",
				Value: ConvertionStatusError,
			},
			{
				Key:   "qc_status",
				Type:  "select",
				Value: ConvertionStatusError,
			},
		}, metadataFields)
		return nil, nil
	}
	c.commentToItem = func(ctx context.Context, assetID, content string) error {
		assert.Contains(t, content, "品質検査・変換に失敗しました。")
		assert.Contains(t, content, "ログ： log")
		return nil
	}
	err = receiveResultFromFME(ctx, s, conf, r)
	assert.NoError(t, err)
}

func getLogs(t *testing.T, f func()) string {
	t.Helper()
	buf := &bytes.Buffer{}
	log.SetOutput(io.MultiWriter(buf, os.Stdout))

	defer func() {
		log.SetOutput(os.Stdout)
	}()

	f()
	return buf.String()
}

type cmsMock struct {
	cms.Interface
	getItem             func(ctx context.Context, id string, asset bool) (*cms.Item, error)
	updateItem          func(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error)
	asset               func(ctx context.Context, id string) (*cms.Asset, error)
	uploadAsset         func(ctx context.Context, projectID, url string) (string, error)
	uploadAssetDirectly func(ctx context.Context, projectID, name string, r io.Reader) (string, error)
	commentToItem       func(ctx context.Context, assetID, content string) error
}

func (c *cmsMock) reset() {
	c.getItem = nil
	c.updateItem = nil
	c.asset = nil
	c.uploadAsset = nil
	c.uploadAssetDirectly = nil
	c.commentToItem = nil
}

func (c *cmsMock) GetItem(ctx context.Context, id string, asset bool) (*cms.Item, error) {
	return c.getItem(ctx, id, asset)
}

func (c *cmsMock) UpdateItem(ctx context.Context, id string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
	return c.updateItem(ctx, id, fields, metadataFields)
}

func (c *cmsMock) Asset(ctx context.Context, id string) (*cms.Asset, error) {
	return c.asset(ctx, id)
}

func (c *cmsMock) UploadAsset(ctx context.Context, projectID, url string) (string, error) {
	return c.uploadAsset(ctx, projectID, url)
}

func (c *cmsMock) UploadAssetDirectly(ctx context.Context, projectID, name string, r io.Reader) (string, error) {
	return c.uploadAssetDirectly(ctx, projectID, name, r)
}

func (c *cmsMock) CommentToItem(ctx context.Context, assetID, content string) error {
	return c.commentToItem(ctx, assetID, content)
}
