package cmsintegrationv3

import (
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

// func tagIs(t *cms.Tag, v  fmt.Stringer) bool {
// 	return t != nil && t.Name == v.String()
// }

func tagIsNot(t *cms.Tag, v fmt.Stringer) bool {
	return t != nil && t.Name != v.String()
}

func tagFrom(t fmt.Stringer) *cms.Tag {
	s := t.String()
	if s == "" {
		return nil
	}
	return &cms.Tag{
		Name: s,
	}
}

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
