package cmsintegrationv2

import (
	cms "github.com/reearth/reearth-cms-api/go"
)

type Item struct {
	ID string `json:"id,omitempty" cms:"id"`
}

func (i Item) Fields() (fields []cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}

func ItemFrom(item cms.Item) (i Item) {
	item.Unmarshal(&i)
	return
}
