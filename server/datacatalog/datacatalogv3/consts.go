package datacatalogv3

var plateauFeatureTypes = []FeatureType{
	{
		Code: "bldg",
		Name: "建築物モデル",
	},
	{
		Code: "tran",
		Name: "交通（道路）モデル",
	},
	{
		Code: "rwy",
		Name: "交通（鉄道）モデル",
	},
	{
		Code: "trk",
		Name: "交通（徒歩道）モデル",
	},
	{
		Code: "squr",
		Name: "交通（広場）モデル",
	},
	{
		Code: "wwy",
		Name: "交通（航路）モデル",
	},
	{
		Code: "luse",
		Name: "土地利用モデル",
	},
	{
		Code: "fld",
		Name: "洪水浸水想定区域モデル",
	},
	{
		Code: "tnm",
		Name: "津波浸水想定区域モデル",
	},
	{
		Code: "htd",
		Name: "高潮浸水想定区域モデル",
	},
	{
		Code: "ifld",
		Name: "内水浸水想定区域モデル",
	},
	{
		Code: "lsld",
		Name: "土砂災害モデル",
	},
	{
		Code: "urf",
		Name: "都市計画決定情報モデル",
	},
	{
		Code: "unf",
		Name: "地下埋設物モデル",
	},
	{
		Code: "brid",
		Name: "橋梁モデル",
	},
	{
		Code: "tun",
		Name: "トンネルモデル",
	},
	{
		Code: "cons",
		Name: "その他の構造物モデル",
	},
	{
		Code: "frn",
		Name: "都市設備モデル",
	},
	{
		Code: "ubld",
		Name: "地下街モデル",
	},
	{
		Code: "veg",
		Name: "植生モデル",
	},
	{
		Code: "dem",
		Name: "地形モデル",
	},
	{
		Code: "wtr",
		Name: "水部モデル",
	},
	{
		Code: "area",
		Name: "区域モデル",
	},
	{
		Code: "gen",
		Name: "汎用都市オブジェクトモデル",
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
		Code: "global",
		Name: "全球データ",
	},
	{
		Code: "usecase",
		Name: "ユースケース",
	},
	{
		Code: "sample",
		Name: "サンプルデータ",
	},
}
