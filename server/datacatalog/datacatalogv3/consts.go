package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

var plateauSpecVersions3 = []string{"3.0", "3.1", "3.2", "3.3"}

var plateauSpecs = []plateauapi.PlateauSpecSimple{
	{
		MajorVersion:  3,
		Year:          2023,
		MinorVersions: plateauSpecVersions3,
	},
}

var plateauFeatureTypes = []FeatureType{
	{
		Code:      "bldg",
		Name:      "建築物モデル",
		SpecMajor: 3,
	},
	{
		Code:      "tran",
		Name:      "交通（道路）モデル",
		SpecMajor: 3,
		MVTLayerNamesForLOD: map[int][]string{
			0: {"Road"},
			1: {"Road"},
			2: {"TrafficArea", "AuxiliaryTrafficArea"},
		},
	},
	{
		Code:         "rwy",
		Name:         "交通（鉄道）モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"rwy"},
	},
	{
		Code:         "trk",
		Name:         "交通（徒歩道）モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"trk"},
	},
	{
		Code:         "squr",
		Name:         "交通（広場）モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"squr"},
	},
	{
		Code:         "wwy",
		Name:         "交通（航路）モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"wwy"},
	},
	{
		Code:         "luse",
		Name:         "土地利用モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"luse"},
	},
	{
		Code:      "fld",
		Name:      "洪水浸水想定区域モデル",
		SpecMajor: 3,
		Flood:     true,
	},
	{
		Code:      "tnm",
		Name:      "津波浸水想定区域モデル",
		SpecMajor: 3,
		Flood:     true,
	},
	{
		Code:      "htd",
		Name:      "高潮浸水想定区域モデル",
		SpecMajor: 3,
		Flood:     true,
	},
	{
		Code:      "ifld",
		Name:      "内水浸水想定区域モデル",
		SpecMajor: 3,
		Flood:     true,
	},
	{
		Code:         "lsld",
		Name:         "土砂災害警戒区域モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"lsld"},
	},
	{
		Code:      "urf",
		Name:      "都市計画決定情報モデル",
		SpecMajor: 3,
	},
	{
		Code:      "unf",
		Name:      "地下埋設物モデル",
		SpecMajor: 3,
	},
	{
		Code:         "brid",
		Name:         "橋梁モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"brid"},
	},
	{
		Code:      "tun",
		Name:      "トンネルモデル",
		SpecMajor: 3,
	},
	{
		Code:      "cons",
		Name:      "その他の構造物モデル",
		SpecMajor: 3,
	},
	{
		Code:      "frn",
		Name:      "都市設備モデル",
		SpecMajor: 3,
	},
	{
		Code:      "ubld",
		Name:      "地下街モデル",
		SpecMajor: 3,
	},
	{
		Code:      "veg",
		Name:      "植生モデル",
		SpecMajor: 3,
	},
	{
		Code:      "dem",
		Name:      "地形モデル",
		SpecMajor: 3,
	},
	{
		Code:         "wtr",
		Name:         "水部モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"wtr"},
	},
	{
		Code:         "area",
		Name:         "区域モデル",
		SpecMajor:    3,
		MVTLayerName: []string{"area"},
	},
	{
		Code:      "gen",
		Name:      "汎用都市オブジェクトモデル",
		SpecMajor: 3,
	},
}

var relatedFeatureTypes = []FeatureType{
	{
		Code: "shelter",
		Name: "避難施設情報",
	},
	{
		Code: "park",
		Name: "公園情報",
	},
	{
		Code: "landmark",
		Name: "ランドマーク情報",
	},
	{
		Code: "station",
		Name: "鉄道駅情報",
	},
	{
		Code: "railway",
		Name: "鉄道情報",
	},
	{
		Code: "emergency_route",
		Name: "緊急輸送道路情報",
	},
	{
		Code: "border",
		Name: "行政界情報",
	},
}

var genericFeatureTypes = []FeatureType{
	{
		Code: "usecase",
		Name: "ユースケース",
	},
	{
		Code: "global",
		Name: "全球データ",
	},
	{
		Code: "sample",
		Name: "サンプルデータ",
	},
}
