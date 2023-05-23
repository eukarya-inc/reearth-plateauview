package datacatalog

import (
	"fmt"

	"golang.org/x/exp/slices"
)

var FeatureTypes = []string{
	"bldg",
	"tran",
	"frn",
	"veg",
	"luse",
	"lsld",
	"urf",
	"fld",
	"htd",
	"ifld",
	"tnm",
	"brid",
	"rail",
	"gen",
}

func (i PlateauItem) DataCatalogItems(c PlateauIntermediateItem, ty string) []*DataCatalogItem {
	// worksround
	switch ty {
	case "fld":
		return i.FldItems(c)
	case "htd":
		return i.HtdItems(c)
	case "ifld":
		return i.IfldItems(c)
	case "tnm":
		return i.TnmItems(c)
	}

	o, ok := FeatureOptions[ty]
	if !ok {
		return nil
	}

	return DataCatalogItemBuilder{
		Assets:           i.Feature(ty),
		Descriptions:     i.FeatureDescription(ty),
		IntermediateItem: c,
		Options:          o,
	}.Build()
}

var FeatureOptions = map[string]DataCatalogItemBuilderOption{
	"bldg": {
		ModelName:          "建築物モデル",
		LOD:                true,
		UseMaxLODAsDefault: true,
		ItemID:             true,
		GroupBy: func(an AssetName) string {
			return an.WardEn
		},
		SortGroupBy: func(a, b string, c, d AssetName) bool {
			return c.WardCodeInt() < d.WardCodeInt() || c.LODInt() < d.LODInt()
		},
	},
	"tran": {
		ModelName: "道路モデル",
		LOD:       true,
		LayersForLOD: map[string][]string{
			"1": {"Road"},
			"2": {"TrafficArea", "AuxiliaryTrafficArea"},
		},
		UseMaxLODAsDefault: true,
	},
	"frn": {
		ModelName: "都市設備モデル",
		LOD:       true,
	},
	"veg": {
		ModelName: "植生モデル",
		LOD:       true,
	},
	"luse": {
		ModelName: "土地利用モデル",
		Layers:    []string{"luse"},
	},
	"lsld": {
		ModelName: "土砂災害警戒区域モデル",
		Layers:    []string{"lsld"},
	},
	"urf": {
		ModelName:           "都市計画決定情報モデル",
		UseGroupNameAsLayer: true,
		MultipleDesc:        true,
		NameOverrideBy: func(an AssetName) (string, string, string) {
			if urfName := urfFeatureTypeMap[an.UrfFeatureType]; urfName != "" {
				return fmt.Sprintf("%sモデル", urfName), urfName, an.UrfFeatureType
			}
			return an.UrfFeatureType, an.UrfFeatureType, an.UrfFeatureType
		},
		GroupBy: func(an AssetName) string {
			return an.UrfFeatureType
		},
		SortGroupBy: func(_, _ string, c, d AssetName) bool {
			i1 := slices.Index(urfFeatureTypes, c.UrfFeatureType)
			if i1 < 0 {
				i1 = len(urfFeatureTypes)
			}
			i2 := slices.Index(urfFeatureTypes, d.UrfFeatureType)
			if i2 < 0 {
				i2 = len(urfFeatureTypes)
			}
			return i1 < i2
		},
	},
	"brid": {
		ModelName: "橋梁モデル",
		LOD:       true,
		Layers:    []string{"brid"},
	},
	"rail": {
		ModelName: "鉄道モデル",
		LOD:       true,
		Layers:    []string{"rail"},
	},
	"gen": {
		ModelName:           "汎用都市オブジェクトモデル",
		MultipleDesc:        true,
		LOD:                 true,
		UseGroupNameAsName:  true,
		UseGroupNameAsLayer: true,
		GroupBy: func(n AssetName) string {
			return n.GenName
		},
	},
}
