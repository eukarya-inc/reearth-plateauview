// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package plateauapi

import (
	"fmt"
	"io"
	"strconv"
)

// 地域。都道府県（Prefecture）・市区町村（City）・区（政令指定都市のみ・Ward）のいずれかです。
// 政令指定都市の場合のみ、市の下に区が存在します。
type Area interface {
	IsNode()
	IsArea()
	GetID() ID
	// 地域の種類
	GetType() AreaType
	// 地域コード。行政コードや市区町村コードとも呼ばれます。
	// 都道府県の場合は二桁の数字から成る文字列です。
	// 市区町村の場合は、先頭に都道府県コードを含む5桁の数字から成る文字列です。
	GetCode() AreaCode
	// 地域名
	GetName() string
	// 地域に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	GetDatasets() []Dataset
}

// データセット。
type Dataset interface {
	IsNode()
	IsDataset()
	GetID() ID
	// データセット名
	GetName() string
	// データセットのサブ名
	GetSubname() *string
	// データセットの説明
	GetDescription() *string
	// データセットの公開年度（西暦）
	GetYear() int
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	GetGroups() []string
	// データセットが属する都道府県のID。
	GetPrefectureID() *ID
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	GetPrefectureCode() *AreaCode
	// データセットが属する市のID。
	GetCityID() *ID
	// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	GetCityCode() *AreaCode
	// データセットが属する区のID。
	GetWardID() *ID
	// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	GetWardCode() *AreaCode
	// データセットの種類のID。
	GetTypeID() ID
	// データセットの種類コード。
	GetTypeCode() string
	// データセットが属する都道府県。
	GetPrefecture() *Prefecture
	// データセットが属する市。
	GetCity() *City
	// データセットが属する区。
	GetWard() *Ward
	// データセットの種類。
	GetType() DatasetType
	// データセットのアイテム。
	GetItems() []DatasetItem
}

// データセットのアイテム。
type DatasetItem interface {
	IsNode()
	IsDatasetItem()
	GetID() ID
	// データセットのアイテムのフォーマット。
	GetFormat() DatasetFormat
	// データセットのアイテム名。
	GetName() string
	// データセットのアイテムのURL。
	GetURL() string
	// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
	// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
	GetLayers() []string
	// データセットのアイテムが属するデータセットのID。
	GetParentID() ID
	// データセットのアイテムが属するデータセット。
	GetParent() Dataset
}

// データセットの種類。
type DatasetType interface {
	IsNode()
	IsDatasetType()
	GetID() ID
	// データセットの種類コード。「bldg」など。
	GetCode() string
	// データセットの種類名。
	GetName() string
	// データセットの種類のカテゴリ。
	GetCategory() DatasetTypeCategory
	// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
	GetDatasets() []Dataset
}

// IDを持つオブジェクト。nodeまたはnodesクエリでIDを指定して検索可能です。
type Node interface {
	IsNode()
	// オブジェクトのID
	GetID() ID
}

// 地域を検索するためのクエリ。
type AreasInput struct {
	// 検索したい地域が属する親となる地域のコード。例えば東京都に属する都市を検索したい場合は "13" を指定します。
	ParentCode *AreaCode `json:"parentCode,omitempty"`
	// データセットの種類コード。例えば、建築物モデルのデータセットが存在する地域を検索したい場合は "bldg" を指定します。複数指定するとOR条件で検索を行います。
	// 未指定の場合、全てのデータセットの種類を対象に検索します。
	DatasetTypes []string `json:"datasetTypes,omitempty"`
	// データセットの種類のカテゴリ。例えば、PLATEAU都市モデルデータセットが存在する地域を検索したい場合は PLATEAU を指定します。複数指定するとOR条件で検索を行います。
	// 未指定の場合、全てのカテゴリのデータセットを対象に検索します。
	Categories []DatasetTypeCategory `json:"categories,omitempty"`
	// 地域の種類。例えば、市を検索したい場合は CITY を指定します。複数指定するとOR条件で検索を行います。
	// 未指定の場合、全ての地域を対象に検索します。
	AreaTypes []AreaType `json:"areaTypes,omitempty"`
	// 検索文字列。複数指定するとAND条件で絞り込み検索が行えます。
	SearchTokens []string `json:"searchTokens,omitempty"`
}

// 市区町村
type City struct {
	ID ID `json:"id"`
	// 地域の種類
	Type AreaType `json:"type"`
	// 市区町村コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 市区町村名
	Name string `json:"name"`
	// 市区町村が属する都道府県のID。
	PrefectureID ID `json:"prefectureId"`
	// 市区町村が属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefectureCode"`
	// 市区町村の都道府県。
	Prefecture *Prefecture `json:"prefecture,omitempty"`
	// 市区町村に属する区。政令指定都市の場合のみ存在します。
	Wards []*Ward `json:"wards"`
	// 市区町村に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (City) IsArea()        {}
func (this City) GetID() ID { return this.ID }

// 地域の種類
func (this City) GetType() AreaType { return this.Type }

// 地域コード。行政コードや市区町村コードとも呼ばれます。
// 都道府県の場合は二桁の数字から成る文字列です。
// 市区町村の場合は、先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this City) GetCode() AreaCode { return this.Code }

// 地域名
func (this City) GetName() string { return this.Name }

// 地域に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
func (this City) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (City) IsNode() {}

// オブジェクトのID

// データセットの種類を検索するためのクエリ。
type DatasetTypesInput struct {
	// データセットの種類のカテゴリ。
	Category *DatasetTypeCategory `json:"category,omitempty"`
	// データセットの種類が属するPLATEAU都市モデルの仕様名。
	PlateauSpec *string `json:"plateauSpec,omitempty"`
	// データセットの種類が属するPLATEAU都市モデルの仕様の公開年度（西暦）。
	Year *int `json:"year,omitempty"`
}

// データセットを検索するためのクエリ。
type DatasetsInput struct {
	// データセットの地域コード（都道府県コードや市区町村コードが使用可能）。複数指定するとOR条件で検索を行います。
	AreaCodes []AreaCode `json:"areaCodes,omitempty"`
	// 仕様書のバージョン。「第2.3版」「2.3」「2」などの文字列が使用可能です。
	PlateauSpec *string `json:"plateauSpec,omitempty"`
	// データの整備年度または公開年度（西暦）。
	Year *int `json:"year,omitempty"`
	// 検索結果から除外するデータセットの種類コード。
	ExcludeTypes []string `json:"excludeTypes,omitempty"`
	// 検索結果に含めるデータセットの種類コード。未指定の場合、全てのデータセットの種類を対象に検索し、指定するとその種類で検索結果を絞り込みます。
	IncludeTypes []string `json:"includeTypes,omitempty"`
	// 検索文字列。複数指定するとAND条件で絞り込み検索が行えます。
	SearchTokens []string `json:"searchTokens,omitempty"`
	// areaCodesで指定された地域に直接属しているデータセットのみを検索対象にするかどうか。
	// デフォルトはfalseで、指定された地域に間接的に属するデータセットも全て検索します。
	// 例えば、札幌市を対象にした場合、札幌市には中央区や北区といった区のデータセットも存在しますが、trueにすると札幌市のデータセットのみを返します。
	Shallow *bool `json:"shallow,omitempty"`
}

// ユースケースデータなどを含む、その他のデータセット。
type GenericDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname,omitempty"`
	// データセットの説明
	Description *string `json:"description,omitempty"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups,omitempty"`
	// データセットが属する都道府県のID。
	PrefectureID *ID `json:"prefectureId,omitempty"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode *AreaCode `json:"prefectureCode,omitempty"`
	// データセットが属する市のID。
	CityID *ID `json:"cityId,omitempty"`
	// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	CityCode *AreaCode `json:"cityCode,omitempty"`
	// データセットが属する区のID。
	WardID *ID `json:"wardId,omitempty"`
	// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	WardCode *AreaCode `json:"wardCode,omitempty"`
	// データセットの種類のID。
	TypeID ID `json:"typeId"`
	// データセットの種類コード。
	TypeCode string `json:"typeCode"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture,omitempty"`
	// データセットが属する市。
	City *City `json:"city,omitempty"`
	// データセットが属する区。
	Ward *Ward `json:"ward,omitempty"`
	// データセットの種類。
	Type *GenericDatasetType `json:"type"`
	// データセットのアイテム。
	Items []*GenericDatasetItem `json:"items"`
}

func (GenericDataset) IsDataset()     {}
func (this GenericDataset) GetID() ID { return this.ID }

// データセット名
func (this GenericDataset) GetName() string { return this.Name }

// データセットのサブ名
func (this GenericDataset) GetSubname() *string { return this.Subname }

// データセットの説明
func (this GenericDataset) GetDescription() *string { return this.Description }

// データセットの公開年度（西暦）
func (this GenericDataset) GetYear() int { return this.Year }

// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
func (this GenericDataset) GetGroups() []string {
	if this.Groups == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Groups))
	for _, concrete := range this.Groups {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットが属する都道府県のID。
func (this GenericDataset) GetPrefectureID() *ID { return this.PrefectureID }

// データセットが属する都道府県コード。2桁の数字から成る文字列です。
func (this GenericDataset) GetPrefectureCode() *AreaCode { return this.PrefectureCode }

// データセットが属する市のID。
func (this GenericDataset) GetCityID() *ID { return this.CityID }

// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this GenericDataset) GetCityCode() *AreaCode { return this.CityCode }

// データセットが属する区のID。
func (this GenericDataset) GetWardID() *ID { return this.WardID }

// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this GenericDataset) GetWardCode() *AreaCode { return this.WardCode }

// データセットの種類のID。
func (this GenericDataset) GetTypeID() ID { return this.TypeID }

// データセットの種類コード。
func (this GenericDataset) GetTypeCode() string { return this.TypeCode }

// データセットが属する都道府県。
func (this GenericDataset) GetPrefecture() *Prefecture { return this.Prefecture }

// データセットが属する市。
func (this GenericDataset) GetCity() *City { return this.City }

// データセットが属する区。
func (this GenericDataset) GetWard() *Ward { return this.Ward }

// データセットの種類。
func (this GenericDataset) GetType() DatasetType { return *this.Type }

// データセットのアイテム。
func (this GenericDataset) GetItems() []DatasetItem {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]DatasetItem, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (GenericDataset) IsNode() {}

// オブジェクトのID

// その他のデータセットのアイテム。
type GenericDatasetItem struct {
	ID ID `json:"id"`
	// データセットのアイテムのフォーマット。
	Format DatasetFormat `json:"format"`
	// データセットのアイテム名。
	Name string `json:"name"`
	// データセットのアイテムのURL。
	URL string `json:"url"`
	// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
	// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
	Layers []string `json:"layers,omitempty"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parentId"`
	// データセットのアイテムが属するデータセット。
	Parent *GenericDataset `json:"parent,omitempty"`
}

func (GenericDatasetItem) IsDatasetItem() {}
func (this GenericDatasetItem) GetID() ID { return this.ID }

// データセットのアイテムのフォーマット。
func (this GenericDatasetItem) GetFormat() DatasetFormat { return this.Format }

// データセットのアイテム名。
func (this GenericDatasetItem) GetName() string { return this.Name }

// データセットのアイテムのURL。
func (this GenericDatasetItem) GetURL() string { return this.URL }

// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
func (this GenericDatasetItem) GetLayers() []string {
	if this.Layers == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Layers))
	for _, concrete := range this.Layers {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットのアイテムが属するデータセットのID。
func (this GenericDatasetItem) GetParentID() ID { return this.ParentID }

// データセットのアイテムが属するデータセット。
func (this GenericDatasetItem) GetParent() Dataset { return *this.Parent }

func (GenericDatasetItem) IsNode() {}

// オブジェクトのID

// その他のデータセットの種類。
type GenericDatasetType struct {
	ID ID `json:"id"`
	// データセットの種類コード。「usecase」など。
	Code string `json:"code"`
	// データセットの種類名。
	Name string `json:"name"`
	// データセットの種類のカテゴリ。
	Category DatasetTypeCategory `json:"category"`
	// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
	Datasets []*GenericDataset `json:"datasets"`
}

func (GenericDatasetType) IsDatasetType() {}
func (this GenericDatasetType) GetID() ID { return this.ID }

// データセットの種類コード。「bldg」など。
func (this GenericDatasetType) GetCode() string { return this.Code }

// データセットの種類名。
func (this GenericDatasetType) GetName() string { return this.Name }

// データセットの種類のカテゴリ。
func (this GenericDatasetType) GetCategory() DatasetTypeCategory { return this.Category }

// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
func (this GenericDatasetType) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (GenericDatasetType) IsNode() {}

// オブジェクトのID

// PLATEAU都市モデルの通常のデータセット。例えば、地物型が建築物モデル（bldg）などのデータセットです。
type PlateauDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname,omitempty"`
	// データセットの説明
	Description *string `json:"description,omitempty"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups,omitempty"`
	// データセットが属する都道府県のID。
	PrefectureID *ID `json:"prefectureId,omitempty"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode *AreaCode `json:"prefectureCode,omitempty"`
	// データセットが属する市のID。
	CityID *ID `json:"cityId,omitempty"`
	// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	CityCode *AreaCode `json:"cityCode,omitempty"`
	// データセットが属する区のID。
	WardID *ID `json:"wardId,omitempty"`
	// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	WardCode *AreaCode `json:"wardCode,omitempty"`
	// データセットの種類のID。
	TypeID ID `json:"typeId"`
	// データセットの種類コード。
	TypeCode string `json:"typeCode"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture,omitempty"`
	// データセットが属する市。
	City *City `json:"city,omitempty"`
	// データセットが属する区。
	Ward *Ward `json:"ward,omitempty"`
	// データセットの種類。
	Type *PlateauDatasetType `json:"type"`
	// データセットのアイテム。
	Items []*PlateauDatasetItem `json:"items"`
	// データセットが準拠するPLATEAU都市モデルの仕様のID。
	PlateauSpecID ID `json:"plateauSpecId"`
	// データセットが準拠するPLATEAU都市モデルの仕様の名称。
	PlateauSpecName string `json:"plateauSpecName"`
	// データセットが準拠するPLATEAU都市モデルの仕様。
	PlateauSpec *PlateauSpecMinor `json:"plateauSpec"`
	// 河川。地物型が洪水浸水想定区域モデル（fld）の場合のみ存在します。
	River *River `json:"river,omitempty"`
}

func (PlateauDataset) IsDataset()     {}
func (this PlateauDataset) GetID() ID { return this.ID }

// データセット名
func (this PlateauDataset) GetName() string { return this.Name }

// データセットのサブ名
func (this PlateauDataset) GetSubname() *string { return this.Subname }

// データセットの説明
func (this PlateauDataset) GetDescription() *string { return this.Description }

// データセットの公開年度（西暦）
func (this PlateauDataset) GetYear() int { return this.Year }

// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
func (this PlateauDataset) GetGroups() []string {
	if this.Groups == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Groups))
	for _, concrete := range this.Groups {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットが属する都道府県のID。
func (this PlateauDataset) GetPrefectureID() *ID { return this.PrefectureID }

// データセットが属する都道府県コード。2桁の数字から成る文字列です。
func (this PlateauDataset) GetPrefectureCode() *AreaCode { return this.PrefectureCode }

// データセットが属する市のID。
func (this PlateauDataset) GetCityID() *ID { return this.CityID }

// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this PlateauDataset) GetCityCode() *AreaCode { return this.CityCode }

// データセットが属する区のID。
func (this PlateauDataset) GetWardID() *ID { return this.WardID }

// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this PlateauDataset) GetWardCode() *AreaCode { return this.WardCode }

// データセットの種類のID。
func (this PlateauDataset) GetTypeID() ID { return this.TypeID }

// データセットの種類コード。
func (this PlateauDataset) GetTypeCode() string { return this.TypeCode }

// データセットが属する都道府県。
func (this PlateauDataset) GetPrefecture() *Prefecture { return this.Prefecture }

// データセットが属する市。
func (this PlateauDataset) GetCity() *City { return this.City }

// データセットが属する区。
func (this PlateauDataset) GetWard() *Ward { return this.Ward }

// データセットの種類。
func (this PlateauDataset) GetType() DatasetType { return *this.Type }

// データセットのアイテム。
func (this PlateauDataset) GetItems() []DatasetItem {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]DatasetItem, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (PlateauDataset) IsNode() {}

// オブジェクトのID

// PLATEAU都市モデルのデータセットのアイテム。
type PlateauDatasetItem struct {
	ID ID `json:"id"`
	// データセットのアイテムのフォーマット。
	Format DatasetFormat `json:"format"`
	// データセットのアイテム名。
	Name string `json:"name"`
	// データセットのアイテムのURL。
	URL string `json:"url"`
	// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
	// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
	Layers []string `json:"layers,omitempty"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parentId"`
	// データセットのアイテムが属するデータセット。
	Parent *PlateauDataset `json:"parent,omitempty"`
	// データセットのアイテムのLOD（詳細度・Level of Detail）。1、2、3、4などの整数値です。
	Lod *int `json:"lod,omitempty"`
	// データセットのアイテムのテクスチャの種類。
	Texture *Texture `json:"texture,omitempty"`
	// 浸水規模。地物型が洪水・高潮・津波・内水浸水想定区域モデル（fld・htd・tnm・ifld）の場合のみ存在します。
	FloodingScale *FloodingScale `json:"floodingScale,omitempty"`
}

func (PlateauDatasetItem) IsDatasetItem() {}
func (this PlateauDatasetItem) GetID() ID { return this.ID }

// データセットのアイテムのフォーマット。
func (this PlateauDatasetItem) GetFormat() DatasetFormat { return this.Format }

// データセットのアイテム名。
func (this PlateauDatasetItem) GetName() string { return this.Name }

// データセットのアイテムのURL。
func (this PlateauDatasetItem) GetURL() string { return this.URL }

// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
func (this PlateauDatasetItem) GetLayers() []string {
	if this.Layers == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Layers))
	for _, concrete := range this.Layers {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットのアイテムが属するデータセットのID。
func (this PlateauDatasetItem) GetParentID() ID { return this.ParentID }

// データセットのアイテムが属するデータセット。
func (this PlateauDatasetItem) GetParent() Dataset { return *this.Parent }

func (PlateauDatasetItem) IsNode() {}

// オブジェクトのID

// PLATEAU都市モデルのデータセットの種類。
type PlateauDatasetType struct {
	ID ID `json:"id"`
	// データセットの種類コード。「bldg」など。
	Code string `json:"code"`
	// データセットの種類名。
	Name string `json:"name"`
	// データセットの種類のカテゴリ。
	Category DatasetTypeCategory `json:"category"`
	// データセットの種類が属するPLATEAU都市モデルの仕様のID。
	PlateauSpecID ID `json:"plateauSpecId"`
	// データセットの種類が属するPLATEAU都市モデルの仕様の名称。
	PlateauSpecName string `json:"plateauSpecName"`
	// データセットの種類が属するPLATEAU都市モデルの仕様。
	PlateauSpec *PlateauSpec `json:"plateauSpec,omitempty"`
	// データセットの種類が属するPLATEAU都市モデルの仕様の公開年度（西暦）。
	Year int `json:"year"`
	// 洪水・高潮・津波・内水浸水想定区域モデルを表す種類かどうか。河川などの情報が利用可能です。
	Flood bool `json:"flood"`
	// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
	Datasets []*PlateauDataset `json:"datasets"`
}

func (PlateauDatasetType) IsDatasetType() {}
func (this PlateauDatasetType) GetID() ID { return this.ID }

// データセットの種類コード。「bldg」など。
func (this PlateauDatasetType) GetCode() string { return this.Code }

// データセットの種類名。
func (this PlateauDatasetType) GetName() string { return this.Name }

// データセットの種類のカテゴリ。
func (this PlateauDatasetType) GetCategory() DatasetTypeCategory { return this.Category }

// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
func (this PlateauDatasetType) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (PlateauDatasetType) IsNode() {}

// オブジェクトのID

// PLATEAU都市モデルの仕様のメジャーバージョン。
type PlateauSpec struct {
	ID ID `json:"id"`
	// PLATEAU都市モデルの仕様のバージョン番号。
	MajorVersion int `json:"majorVersion"`
	// 仕様の公開年度（西暦）。
	Year int `json:"year"`
	// その仕様に含まれるデータセットの種類。
	DatasetTypes []*PlateauDatasetType `json:"datasetTypes"`
	// その仕様のマイナーバージョン。
	MinorVersions []*PlateauSpecMinor `json:"minorVersions"`
}

func (PlateauSpec) IsNode() {}

// オブジェクトのID
func (this PlateauSpec) GetID() ID { return this.ID }

// PLATEAU都市モデルの仕様のマイナーバージョン。
type PlateauSpecMinor struct {
	ID ID `json:"id"`
	// PLATEAU都市モデルの仕様の名前。 "第2.3版" のような文字列です。
	Name string `json:"name"`
	// バージョンを表す文字列。 "2.3" のような文字列です。
	Version string `json:"version"`
	// メジャーバージョン番号。 2のような整数です。
	MajorVersion int `json:"majorVersion"`
	// 仕様の公開年度（西暦）。
	Year int `json:"year"`
	// その仕様が属する仕様のメジャーバージョンのID。
	ParentID ID `json:"parentId"`
	// その仕様が属する仕様のメジャーバージョン。
	Parent *PlateauSpec `json:"parent"`
	// その仕様に準拠して整備されたPLATEAU都市モデルデータセット（DatasetInput内のplateauSpecの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (PlateauSpecMinor) IsNode() {}

// オブジェクトのID
func (this PlateauSpecMinor) GetID() ID { return this.ID }

// 都道府県
type Prefecture struct {
	ID ID `json:"id"`
	// 地域の種類
	Type AreaType `json:"type"`
	// 都道府県コード。2桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 都道府県名
	Name string `json:"name"`
	// 都道府県に属する市区町村
	Cities []*City `json:"cities"`
	// 都道府県に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (Prefecture) IsArea()        {}
func (this Prefecture) GetID() ID { return this.ID }

// 地域の種類
func (this Prefecture) GetType() AreaType { return this.Type }

// 地域コード。行政コードや市区町村コードとも呼ばれます。
// 都道府県の場合は二桁の数字から成る文字列です。
// 市区町村の場合は、先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this Prefecture) GetCode() AreaCode { return this.Code }

// 地域名
func (this Prefecture) GetName() string { return this.Name }

// 地域に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
func (this Prefecture) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (Prefecture) IsNode() {}

// オブジェクトのID

// PLATEAU都市モデルデータセットと併せて表示することで情報を補完できる、関連データセット。
// 避難施設・ランドマーク・鉄道駅・鉄道・緊急輸送道路・公園・行政界などのデータセット。
type RelatedDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname,omitempty"`
	// データセットの説明
	Description *string `json:"description,omitempty"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups,omitempty"`
	// データセットが属する都道府県のID。
	PrefectureID *ID `json:"prefectureId,omitempty"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode *AreaCode `json:"prefectureCode,omitempty"`
	// データセットが属する市のID。
	CityID *ID `json:"cityId,omitempty"`
	// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	CityCode *AreaCode `json:"cityCode,omitempty"`
	// データセットが属する区のID。
	WardID *ID `json:"wardId,omitempty"`
	// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	WardCode *AreaCode `json:"wardCode,omitempty"`
	// データセットの種類のID。
	TypeID ID `json:"typeId"`
	// データセットの種類コード。
	TypeCode string `json:"typeCode"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture,omitempty"`
	// データセットが属する市。
	City *City `json:"city,omitempty"`
	// データセットが属する区。
	Ward *Ward `json:"ward,omitempty"`
	// データセットの種類。
	Type *RelatedDatasetType `json:"type"`
	// データセットのアイテム。
	Items []*RelatedDatasetItem `json:"items"`
}

func (RelatedDataset) IsDataset()     {}
func (this RelatedDataset) GetID() ID { return this.ID }

// データセット名
func (this RelatedDataset) GetName() string { return this.Name }

// データセットのサブ名
func (this RelatedDataset) GetSubname() *string { return this.Subname }

// データセットの説明
func (this RelatedDataset) GetDescription() *string { return this.Description }

// データセットの公開年度（西暦）
func (this RelatedDataset) GetYear() int { return this.Year }

// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
func (this RelatedDataset) GetGroups() []string {
	if this.Groups == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Groups))
	for _, concrete := range this.Groups {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットが属する都道府県のID。
func (this RelatedDataset) GetPrefectureID() *ID { return this.PrefectureID }

// データセットが属する都道府県コード。2桁の数字から成る文字列です。
func (this RelatedDataset) GetPrefectureCode() *AreaCode { return this.PrefectureCode }

// データセットが属する市のID。
func (this RelatedDataset) GetCityID() *ID { return this.CityID }

// データセットが属する市コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this RelatedDataset) GetCityCode() *AreaCode { return this.CityCode }

// データセットが属する区のID。
func (this RelatedDataset) GetWardID() *ID { return this.WardID }

// データセットが属する区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this RelatedDataset) GetWardCode() *AreaCode { return this.WardCode }

// データセットの種類のID。
func (this RelatedDataset) GetTypeID() ID { return this.TypeID }

// データセットの種類コード。
func (this RelatedDataset) GetTypeCode() string { return this.TypeCode }

// データセットが属する都道府県。
func (this RelatedDataset) GetPrefecture() *Prefecture { return this.Prefecture }

// データセットが属する市。
func (this RelatedDataset) GetCity() *City { return this.City }

// データセットが属する区。
func (this RelatedDataset) GetWard() *Ward { return this.Ward }

// データセットの種類。
func (this RelatedDataset) GetType() DatasetType { return *this.Type }

// データセットのアイテム。
func (this RelatedDataset) GetItems() []DatasetItem {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]DatasetItem, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (RelatedDataset) IsNode() {}

// オブジェクトのID

// 関連データセットのアイテム。
type RelatedDatasetItem struct {
	ID ID `json:"id"`
	// データセットのアイテムのフォーマット。
	Format DatasetFormat `json:"format"`
	// データセットのアイテム名。
	Name string `json:"name"`
	// データセットのアイテムのURL。
	URL string `json:"url"`
	// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
	// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
	Layers []string `json:"layers,omitempty"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parentId"`
	// データセットのアイテムが属するデータセット。
	Parent *RelatedDataset `json:"parent,omitempty"`
}

func (RelatedDatasetItem) IsDatasetItem() {}
func (this RelatedDatasetItem) GetID() ID { return this.ID }

// データセットのアイテムのフォーマット。
func (this RelatedDatasetItem) GetFormat() DatasetFormat { return this.Format }

// データセットのアイテム名。
func (this RelatedDatasetItem) GetName() string { return this.Name }

// データセットのアイテムのURL。
func (this RelatedDatasetItem) GetURL() string { return this.URL }

// データセットのアイテムのレイヤー名。MVTやWMSなどのフォーマットの場合のみ存在。
// レイヤー名が複数存在する場合は、同時に複数のレイヤーを表示可能であることを意味します。
func (this RelatedDatasetItem) GetLayers() []string {
	if this.Layers == nil {
		return nil
	}
	interfaceSlice := make([]string, 0, len(this.Layers))
	for _, concrete := range this.Layers {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

// データセットのアイテムが属するデータセットのID。
func (this RelatedDatasetItem) GetParentID() ID { return this.ParentID }

// データセットのアイテムが属するデータセット。
func (this RelatedDatasetItem) GetParent() Dataset { return *this.Parent }

func (RelatedDatasetItem) IsNode() {}

// オブジェクトのID

// 関連データセットの種類。
type RelatedDatasetType struct {
	ID ID `json:"id"`
	// データセットの種類コード。「park」など。
	Code string `json:"code"`
	// データセットの種類名。
	Name string `json:"name"`
	// データセットの種類のカテゴリ。
	Category DatasetTypeCategory `json:"category"`
	// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
	Datasets []*RelatedDataset `json:"datasets"`
}

func (RelatedDatasetType) IsDatasetType() {}
func (this RelatedDatasetType) GetID() ID { return this.ID }

// データセットの種類コード。「bldg」など。
func (this RelatedDatasetType) GetCode() string { return this.Code }

// データセットの種類名。
func (this RelatedDatasetType) GetName() string { return this.Name }

// データセットの種類のカテゴリ。
func (this RelatedDatasetType) GetCategory() DatasetTypeCategory { return this.Category }

// データセット（DatasetInput内のincludeTypesとexcludeTypesの指定は無視されます）。
func (this RelatedDatasetType) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (RelatedDatasetType) IsNode() {}

// オブジェクトのID

// 洪水浸水想定区域モデルにおける河川。
type River struct {
	// 河川名。通常、「〜水系〜川」という形式になります。
	Name string `json:"name"`
	// 管理区間
	Admin RiverAdmin `json:"admin"`
}

// 区（政令指定都市のみ）
type Ward struct {
	ID ID `json:"id"`
	// 種類
	Type AreaType `json:"type"`
	// 区コード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 区名
	Name string `json:"name"`
	// 区が属する都道府県のID。
	PrefectureID ID `json:"prefectureId"`
	// 区が属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefectureCode"`
	// 区が属する市のID。
	CityID ID `json:"cityId"`
	// 区が属する市のコード。先頭に都道府県コードを含む5桁の数字から成る文字列です。
	CityCode AreaCode `json:"cityCode"`
	// 区が属する都道府県。
	Prefecture *Prefecture `json:"prefecture,omitempty"`
	// 区が属する市。
	City *City `json:"city,omitempty"`
	// 区に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (Ward) IsArea()        {}
func (this Ward) GetID() ID { return this.ID }

// 地域の種類
func (this Ward) GetType() AreaType { return this.Type }

// 地域コード。行政コードや市区町村コードとも呼ばれます。
// 都道府県の場合は二桁の数字から成る文字列です。
// 市区町村の場合は、先頭に都道府県コードを含む5桁の数字から成る文字列です。
func (this Ward) GetCode() AreaCode { return this.Code }

// 地域名
func (this Ward) GetName() string { return this.Name }

// 地域に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
func (this Ward) GetDatasets() []Dataset {
	if this.Datasets == nil {
		return nil
	}
	interfaceSlice := make([]Dataset, 0, len(this.Datasets))
	for _, concrete := range this.Datasets {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}

func (Ward) IsNode() {}

// オブジェクトのID

type AreaType string

const (
	// 都道府県
	AreaTypePrefecture AreaType = "PREFECTURE"
	// 市町村
	AreaTypeCity AreaType = "CITY"
	// 区（政令指定都市のみ）
	AreaTypeWard AreaType = "WARD"
)

var AllAreaType = []AreaType{
	AreaTypePrefecture,
	AreaTypeCity,
	AreaTypeWard,
}

func (e AreaType) IsValid() bool {
	switch e {
	case AreaTypePrefecture, AreaTypeCity, AreaTypeWard:
		return true
	}
	return false
}

func (e AreaType) String() string {
	return string(e)
}

func (e *AreaType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AreaType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AreaType", str)
	}
	return nil
}

func (e AreaType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// データセットのフォーマット。
type DatasetFormat string

const (
	// CSV
	DatasetFormatCSV DatasetFormat = "CSV"
	// CZML
	DatasetFormatCzml DatasetFormat = "CZML"
	// 3D Tiles
	DatasetFormatCesium3dtiles DatasetFormat = "CESIUM3DTILES"
	// GlTF
	DatasetFormatGltf DatasetFormat = "GLTF"
	// GTFS Realtime
	DatasetFormatGtfsRealtime DatasetFormat = "GTFS_REALTIME"
	// GeoJSON
	DatasetFormatGeojson DatasetFormat = "GEOJSON"
	// Mapbox Vector Tile
	DatasetFormatMvt DatasetFormat = "MVT"
	// Tile Map Service
	DatasetFormatTms DatasetFormat = "TMS"
	// XYZで分割された画像タイル。/{z}/{x}/{y}.png のようなURLになります。
	DatasetFormatTiles DatasetFormat = "TILES"
	// Web Map Service
	DatasetFormatWms DatasetFormat = "WMS"
)

var AllDatasetFormat = []DatasetFormat{
	DatasetFormatCSV,
	DatasetFormatCzml,
	DatasetFormatCesium3dtiles,
	DatasetFormatGltf,
	DatasetFormatGtfsRealtime,
	DatasetFormatGeojson,
	DatasetFormatMvt,
	DatasetFormatTms,
	DatasetFormatTiles,
	DatasetFormatWms,
}

func (e DatasetFormat) IsValid() bool {
	switch e {
	case DatasetFormatCSV, DatasetFormatCzml, DatasetFormatCesium3dtiles, DatasetFormatGltf, DatasetFormatGtfsRealtime, DatasetFormatGeojson, DatasetFormatMvt, DatasetFormatTms, DatasetFormatTiles, DatasetFormatWms:
		return true
	}
	return false
}

func (e DatasetFormat) String() string {
	return string(e)
}

func (e *DatasetFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DatasetFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DatasetFormat", str)
	}
	return nil
}

func (e DatasetFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// データセットの種類のカテゴリ。
type DatasetTypeCategory string

const (
	// PLATEAU都市モデルデータセット
	DatasetTypeCategoryPlateau DatasetTypeCategory = "PLATEAU"
	// 関連データセット
	DatasetTypeCategoryRelated DatasetTypeCategory = "RELATED"
	// その他のデータセット
	DatasetTypeCategoryGeneric DatasetTypeCategory = "GENERIC"
)

var AllDatasetTypeCategory = []DatasetTypeCategory{
	DatasetTypeCategoryPlateau,
	DatasetTypeCategoryRelated,
	DatasetTypeCategoryGeneric,
}

func (e DatasetTypeCategory) IsValid() bool {
	switch e {
	case DatasetTypeCategoryPlateau, DatasetTypeCategoryRelated, DatasetTypeCategoryGeneric:
		return true
	}
	return false
}

func (e DatasetTypeCategory) String() string {
	return string(e)
}

func (e *DatasetTypeCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DatasetTypeCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DatasetTypeCategory", str)
	}
	return nil
}

func (e DatasetTypeCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// 浸水想定区域モデルにおける浸水規模。
type FloodingScale string

const (
	// 計画規模
	FloodingScalePlanned FloodingScale = "PLANNED"
	// 想定最大規模
	FloodingScaleExpectedMaximum FloodingScale = "EXPECTED_MAXIMUM"
)

var AllFloodingScale = []FloodingScale{
	FloodingScalePlanned,
	FloodingScaleExpectedMaximum,
}

func (e FloodingScale) IsValid() bool {
	switch e {
	case FloodingScalePlanned, FloodingScaleExpectedMaximum:
		return true
	}
	return false
}

func (e FloodingScale) String() string {
	return string(e)
}

func (e *FloodingScale) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FloodingScale(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FloodingScale", str)
	}
	return nil
}

func (e FloodingScale) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// 河川の管理区間
type RiverAdmin string

const (
	// 国管理区間
	RiverAdminNational RiverAdmin = "NATIONAL"
	// 都道府県管理区間
	RiverAdminPrefecture RiverAdmin = "PREFECTURE"
)

var AllRiverAdmin = []RiverAdmin{
	RiverAdminNational,
	RiverAdminPrefecture,
}

func (e RiverAdmin) IsValid() bool {
	switch e {
	case RiverAdminNational, RiverAdminPrefecture:
		return true
	}
	return false
}

func (e RiverAdmin) String() string {
	return string(e)
}

func (e *RiverAdmin) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RiverAdmin(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RiverAdmin", str)
	}
	return nil
}

func (e RiverAdmin) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// 建築物モデルのテクスチャの種類。
type Texture string

const (
	// テクスチャなし
	TextureNone Texture = "NONE"
	// テクスチャあり
	TextureTexture Texture = "TEXTURE"
)

var AllTexture = []Texture{
	TextureNone,
	TextureTexture,
}

func (e Texture) IsValid() bool {
	switch e {
	case TextureNone, TextureTexture:
		return true
	}
	return false
}

func (e Texture) String() string {
	return string(e)
}

func (e *Texture) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Texture(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Texture", str)
	}
	return nil
}

func (e Texture) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
