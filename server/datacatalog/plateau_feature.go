package datacatalog

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
	case "bldg":
		return i.BldgItems(c)
	case "tran":
		return []*DataCatalogItem{i.TranItem(c)}
	case "luse":
		return []*DataCatalogItem{i.LuseItem(c)}
	case "lsld":
		return []*DataCatalogItem{i.LsldItem(c)}
	case "urf":
		return i.UrfItems(c)
	case "fld":
		return i.FldItems(c)
	case "htd":
		return i.HtdItems(c)
	case "ifld":
		return i.IfldItems(c)
	case "tnm":
		return i.TnmItems(c)
	case "gen":
		return i.GenItems(c)
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
	"frn": {
		ModelName: "都市設備モデル",
		LOD:       true,
	},
	"veg": {
		ModelName: "植生モデル",
		LOD:       true,
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
}
