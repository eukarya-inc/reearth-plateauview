package indexer

import (
	"github.com/qmuntal/gltf"
	"github.com/reearth/go3dtiles/b3dm"

	"errors"
)

func GetGltfAttribute(primitive *gltf.Primitive, doc *gltf.Document, name string) ([]interface{}, error) {
	accessors := doc.Accessors
	attributes := primitive.Attributes
	if len(attributes) == 0 {
		return nil, errors.New("no attributes found")
	}
	att, ok := attributes[name]
	if !ok {
		return nil, errors.New("can't access attribute")
	}
	count := accessors[att].Count

	var res []interface{}
	for i := uint32(0); i < count; i++ {
		res = append(res, b3dm.ReadGltfValueAt(doc, att, i))
	}

	return res, nil
}
