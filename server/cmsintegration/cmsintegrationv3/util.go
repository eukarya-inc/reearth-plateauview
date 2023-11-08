package cmsintegrationv3

import (
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/samber/lo"
)

func updatedField(w *cmswebhook.Payload, field string) *cms.Value {
	for _, c := range w.ItemData.Changes {
		if c.Type == "" || c.Type == cms.FieldChangeTypeDelete {
			continue
		}

		f, _ := lo.Find(w.ItemData.Item.Fields, func(f *cms.Field) bool {
			return f.ID == c.ID
		})
		if f == nil || f.Key != field {
			continue
		}

		return c.GetCurrentValue()
	}

	return nil
}
