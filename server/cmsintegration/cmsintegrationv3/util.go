package cmsintegrationv3

import (
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

func getFieldChangeByKey(w *cmswebhook.ItemData, key string) *cms.FieldChange {
	f := w.Item.FieldByKey(key)
	if f == nil {
		f = w.Item.MetadataField(key)
	}
	if f == nil {
		return nil
	}

	for _, c := range w.Changes {
		if c.ID == f.ID {
			return &c
		}
	}

	return nil
}
