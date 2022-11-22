package indexer

type IndexRoot struct {
	ResultDataUrl string
	IdProperty    string
	Indexes       map[string]interface{}
}

type EnumIndex struct {
	Kind   string
	Values map[string]*EnumValue
}

type EnumValue struct {
	Count      int
	Url        string
	DataRowIds []int
}
