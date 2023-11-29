package cmsintegrationv3

import (
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/samber/lo"
)

func FindFieldChangeByKey(w *cmswebhook.Payload, fieldKey string) (any, bool) {
	if w == nil || w.ItemData == nil || w.ItemData.Item == nil {
		return nil, false
	}

	f := w.ItemData.Item.FieldByKey(fieldKey)
	if f == nil {
		f = w.ItemData.Item.MetadataFieldByKey(fieldKey)
		if f == nil {
			return nil, false
		}
	}

	c, ok := lo.Find(w.ItemData.Changes, func(c cms.FieldChange) bool {
		return c.ID == f.ID
	})
	if !ok {
		return nil, false
	}

	return c.CurrentValue, true
}
