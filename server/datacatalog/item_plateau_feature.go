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

	return DataCatalogItemBuilder{
		Assets:           i.Feature(ty),
		Description:      i.DescriptionVeg,
		IntermediateItem: c,
		Options:          FeatureOptions[ty],
	}.Build()
}

var FeatureOptions = map[string]DataCatalogItemBuilderOption{
	"frn":  FrnOption,
	"veg":  VegOption,
	"brid": BridOption,
	"rail": RailOption,
}

var FrnOption = DataCatalogItemBuilderOption{
	ModelName:   "都市設備モデル",
	MultipleLOD: true,
}

var VegOption = DataCatalogItemBuilderOption{
	ModelName:   "植生モデル",
	MultipleLOD: true,
}

var BridOption = DataCatalogItemBuilderOption{
	ModelName:   "橋梁モデル",
	MultipleLOD: true,
	Layers:      []string{"brid"},
}

var RailOption = DataCatalogItemBuilderOption{
	ModelName:   "鉄道モデル",
	MultipleLOD: true,
	Layers:      []string{"rail"},
}
