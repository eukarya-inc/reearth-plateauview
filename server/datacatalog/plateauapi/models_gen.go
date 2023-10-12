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
	Node
	IsArea()
}

// データセット。
type Dataset interface {
	Node
	IsDataset()
}

// データセットのアイテム。
type DatasetItem interface {
	Node
	IsDatasetItem()
}

// データセットの種類。
type DatasetType interface {
	Node
	IsDatasetType()
}

// IDを持つオブジェクト。nodeまたはnodesクエリでIDを指定して検索可能です。
type Node interface {
	IsNode()
}

// 地域を検索するためのクエリ。
type AreaInput struct {
	// 検索したい地域が属する親となる地域のコード。例えば東京都に属する都市を検索したい場合は "13" を指定します。
	ParentCode *AreaCode `json:"parentCode"`
	// データセットの種類コード。例えば、建築物モデルのデータセットが存在する地域を検索したい場合は "bldg" を指定します。複数指定するとOR条件で検索を行います。
	DatasetTypes []string `json:"datasetTypes"`
	// 検索文字列。複数指定するとAND条件で絞り込み検索が行えます。
	SearchTokens []string `json:"searchTokens"`
}

// 市区町村
type City struct {
	ID ID `json:"id"`
	// 市区町村コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 市区町村名
	Name string `json:"name"`
	// 市区町村が属する都道府県のID。
	PrefectureID ID `json:"prefecture_id"`
	// 市区町村が属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefecture_code"`
	// 市区町村の都道府県。
	Prefecture *Prefecture `json:"prefecture"`
	// 市区町村に属する区。政令指定都市の場合のみ存在します。
	Wards []*Ward `json:"wards"`
	// 市区町村に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (City) IsArea() {}
func (City) IsNode() {}

// データセットを検索するためのクエリ。
type DatasetInput struct {
	// データセットの地域コード（都道府県コードや市区町村コードが使用可能）。複数指定するとOR条件で検索を行います。
	AreaCodes []AreaCode `json:"areaCodes"`
	// 仕様書のバージョン。「第2.3版」「2.3」「2」などの文字列が使用可能です。
	PlateauSpec *string `json:"plateauSpec"`
	// データの整備年度または公開年度（西暦）。
	Year *int `json:"year"`
	// 検索結果から除外するデータセットの種類コード。
	ExcludeTypes []string `json:"excludeTypes"`
	// 検索結果に含めるデータセットの種類コード。未指定の場合、全てのデータセットの種類を対象に検索し、指定するとその種類で検索結果を絞り込みます。
	IncludeTypes []string `json:"includeTypes"`
	// 検索文字列。複数指定するとAND条件で絞り込み検索が行えます。
	SearchTokens []string `json:"searchTokens"`
	// areaCodesで指定された地域に直接属しているデータセットのみを検索対象にするかどうか。
	// デフォルトはfalseで、指定された地域に間接的に属するデータセットも全て検索します。
	// 例えば、札幌市を対象にした場合、札幌市には中央区や北区といった区のデータセットも存在しますが、trueにすると札幌市のデータセットのみを返します。
	Shallow *bool `json:"shallow"`
}

// データセットの種類を検索するためのクエリ。
type DatasetTypeInput struct {
	// データセットの種類のカテゴリ。
	Category *DatasetTypeCategory `json:"category"`
	// データセットの種類が属するPLATEAU都市モデルの仕様名。
	PlateauSpec *string `json:"plateauSpec"`
	// データセットの種類が属するPLATEAU都市モデルの仕様の公開年度（西暦）。
	Year *int `json:"year"`
}

// ユースケースデータなどを含む、その他のデータセット。
type GenericDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname"`
	// データセットの説明
	Description *string `json:"description"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups"`
	// データセットが属する都道府県のID。
	PrefectureID ID `json:"prefecture_id"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefecture_code"`
	// データセットが属する市のID。
	CityID *ID `json:"city_id"`
	// データセットが属する市コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	CityCode *AreaCode `json:"city_code"`
	// データセットが属する区のID。
	WardID *ID `json:"ward_id"`
	// データセットが属する区コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	WardCode *AreaCode `json:"ward_code"`
	// データセットの種類のID。
	TypeID ID `json:"type_id"`
	// データセットの種類コード。
	TypeCode string `json:"type_code"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture"`
	// データセットが属する市。
	City *City `json:"city"`
	// データセットが属する区。
	Ward *Ward `json:"ward"`
	// データセットの種類。
	Type *GenericDatasetType `json:"type"`
	// データセットのアイテム。
	Items []*GenericDatasetItem `json:"items"`
}

func (GenericDataset) IsDataset() {}
func (GenericDataset) IsNode()    {}

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
	Layers []string `json:"layers"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parent_id"`
	// データセットのアイテムが属するデータセット。
	Parent *GenericDataset `json:"parent"`
}

func (GenericDatasetItem) IsDatasetItem() {}
func (GenericDatasetItem) IsNode()        {}

// その他のデータセットの種類。
type GenericDatasetType struct {
	ID ID `json:"id"`
	// データセットの種類コード。「usecase」など。
	Code string `json:"code"`
	// データセットの種類名。
	Name string `json:"name"`
	// データセットの種類のカテゴリ。
	Category DatasetTypeCategory `json:"category"`
}

func (GenericDatasetType) IsDatasetType() {}
func (GenericDatasetType) IsNode()        {}

// PLATEAU都市モデルの通常のデータセット。例えば、地物型が建築物モデル（bldg）などのデータセットです。
type PlateauDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname"`
	// データセットの説明
	Description *string `json:"description"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups"`
	// データセットが属する都道府県のID。
	PrefectureID ID `json:"prefecture_id"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefecture_code"`
	// データセットが属する市のID。
	CityID *ID `json:"city_id"`
	// データセットが属する市コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	CityCode *AreaCode `json:"city_code"`
	// データセットが属する区のID。
	WardID *ID `json:"ward_id"`
	// データセットが属する区コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	WardCode *AreaCode `json:"ward_code"`
	// データセットの種類のID。
	TypeID ID `json:"type_id"`
	// データセットの種類コード。
	TypeCode string `json:"type_code"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture"`
	// データセットが属する市。
	City *City `json:"city"`
	// データセットが属する区。
	Ward *Ward `json:"ward"`
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
	River *River `json:"river"`
}

func (PlateauDataset) IsDataset() {}
func (PlateauDataset) IsNode()    {}

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
	Layers []string `json:"layers"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parent_id"`
	// データセットのアイテムが属するデータセット。
	Parent *PlateauDataset `json:"parent"`
	// データセットのアイテムのLOD（詳細度・Level of Detail）。1、2、3、4などの整数値です。
	Lod *int `json:"lod"`
	// データセットのアイテムのテクスチャの種類。
	Texture *Texture `json:"texture"`
	// 浸水規模。地物型が洪水・高潮・津波・内水浸水想定区域モデル（fld・htd・tnm・ifld）の場合のみ存在します。
	FloodingScale *FloodingScale `json:"floodingScale"`
}

func (PlateauDatasetItem) IsDatasetItem() {}
func (PlateauDatasetItem) IsNode()        {}

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
	PlateauSpec *PlateauSpec `json:"plateauSpec"`
	// データセットの種類が属するPLATEAU都市モデルの仕様の公開年度（西暦）。
	Year int `json:"year"`
	// 洪水・高潮・津波・内水浸水想定区域モデルを表す種類かどうか。河川などの情報が利用可能です。
	Flood bool `json:"flood"`
}

func (PlateauDatasetType) IsDatasetType() {}
func (PlateauDatasetType) IsNode()        {}

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

// 都道府県
type Prefecture struct {
	ID ID `json:"id"`
	// 都道府県コード。2桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 都道府県名
	Name string `json:"name"`
	// 都道府県に属する市区町村
	Cities []*City `json:"cities"`
	// 都道府県に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (Prefecture) IsArea() {}
func (Prefecture) IsNode() {}

// PLATEAU都市モデルデータセットと併せて表示することで情報を補完できる、関連データセット。
// 避難施設・ランドマーク・鉄道駅・鉄道・緊急輸送道路・公園・行政界などのデータセット。
type RelatedDataset struct {
	ID ID `json:"id"`
	// データセット名
	Name string `json:"name"`
	// データセットのサブ名
	Subname *string `json:"subname"`
	// データセットの説明
	Description *string `json:"description"`
	// データセットの公開年度（西暦）
	Year int `json:"year"`
	// データセットを分類するグループ。グループが階層構造になっている場合は、親から子の順番で複数のグループ名が存在することがあります。
	Groups []string `json:"groups"`
	// データセットが属する都道府県のID。
	PrefectureID ID `json:"prefecture_id"`
	// データセットが属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefecture_code"`
	// データセットが属する市のID。
	CityID *ID `json:"city_id"`
	// データセットが属する市コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	CityCode *AreaCode `json:"city_code"`
	// データセットが属する区のID。
	WardID *ID `json:"ward_id"`
	// データセットが属する区コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	WardCode *AreaCode `json:"ward_code"`
	// データセットの種類のID。
	TypeID ID `json:"type_id"`
	// データセットの種類コード。
	TypeCode string `json:"type_code"`
	// データセットが属する都道府県。
	Prefecture *Prefecture `json:"prefecture"`
	// データセットが属する市。
	City *City `json:"city"`
	// データセットが属する区。
	Ward *Ward `json:"ward"`
	// データセットの種類。
	Type *RelatedDatasetType `json:"type"`
	// データセットのアイテム。
	Items []*RelatedDatasetItem `json:"items"`
}

func (RelatedDataset) IsDataset() {}
func (RelatedDataset) IsNode()    {}

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
	Layers []string `json:"layers"`
	// データセットのアイテムが属するデータセットのID。
	ParentID ID `json:"parent_id"`
	// データセットのアイテムが属するデータセット。
	Parent *RelatedDataset `json:"parent"`
}

func (RelatedDatasetItem) IsDatasetItem() {}
func (RelatedDatasetItem) IsNode()        {}

// 関連データセットの種類。
type RelatedDatasetType struct {
	ID ID `json:"id"`
	// データセットの種類コード。「park」など。
	Code string `json:"code"`
	// データセットの種類名。
	Name string `json:"name"`
	// データセットの種類のカテゴリ。
	Category DatasetTypeCategory `json:"category"`
}

func (RelatedDatasetType) IsDatasetType() {}
func (RelatedDatasetType) IsNode()        {}

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
	// 区コード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	Code AreaCode `json:"code"`
	// 区名
	Name string `json:"name"`
	// 区が属する都道府県のID。
	PrefectureID ID `json:"prefecture_id"`
	// 区が属する都道府県コード。2桁の数字から成る文字列です。
	PrefectureCode AreaCode `json:"prefecture_code"`
	// 区が属する市のID。
	CityID ID `json:"city_id"`
	// 区が属する市のコード。先頭に都道府県コードを含む6桁の数字から成る文字列です。
	CityCode AreaCode `json:"city_code"`
	// 区が属する都道府県。
	Prefecture *Prefecture `json:"prefecture"`
	// 区が属する市。
	City *City `json:"city"`
	// 区に属するデータセット（DatasetInput内のareasCodeの指定は無視されます）。
	Datasets []Dataset `json:"datasets"`
}

func (Ward) IsArea() {}
func (Ward) IsNode() {}

// データセットのフォーマット。
type DatasetFormat string

const (
	// CSV
	DatasetFormatCSV DatasetFormat = "CSV"
	// CZML
	DatasetFormatCzml DatasetFormat = "CZML"
	// 3D Tiles
	DatasetFormatCesium3DTiles DatasetFormat = "Cesium3DTiles"
	// GlTF
	DatasetFormatGltf DatasetFormat = "GLTF"
	// GTFS Realtime
	DatasetFormatGTFSRelatime DatasetFormat = "GTFSRelatime"
	// GeoJSON
	DatasetFormatGeoJSON DatasetFormat = "GeoJSON"
	// Mapbox Vector Tile
	DatasetFormatMvt DatasetFormat = "MVT"
	// Tile Map Service
	DatasetFormatTms DatasetFormat = "TMS"
	// XYZで分割された画像タイル。/{z}/{x}/{y}.png のようなURLになります。
	DatasetFormatTiles DatasetFormat = "Tiles"
	// Web Map Service
	DatasetFormatWms DatasetFormat = "WMS"
)

var AllDatasetFormat = []DatasetFormat{
	DatasetFormatCSV,
	DatasetFormatCzml,
	DatasetFormatCesium3DTiles,
	DatasetFormatGltf,
	DatasetFormatGTFSRelatime,
	DatasetFormatGeoJSON,
	DatasetFormatMvt,
	DatasetFormatTms,
	DatasetFormatTiles,
	DatasetFormatWms,
}

func (e DatasetFormat) IsValid() bool {
	switch e {
	case DatasetFormatCSV, DatasetFormatCzml, DatasetFormatCesium3DTiles, DatasetFormatGltf, DatasetFormatGTFSRelatime, DatasetFormatGeoJSON, DatasetFormatMvt, DatasetFormatTms, DatasetFormatTiles, DatasetFormatWms:
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
	DatasetTypeCategoryPlateau DatasetTypeCategory = "Plateau"
	// 関連データセット
	DatasetTypeCategoryRelated DatasetTypeCategory = "Related"
	// その他のデータセット
	DatasetTypeCategoryGeneric DatasetTypeCategory = "Generic"
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
