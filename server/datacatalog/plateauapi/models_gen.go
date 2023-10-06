// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package plateauapi

import (
	"fmt"
	"io"
	"strconv"
)

type Area interface {
	Node
	IsArea()
}

type Dataset interface {
	Node
	IsDataset()
}

type DatasetItem interface {
	Node
	IsDatasetItem()
}

type DatasetType interface {
	Node
	IsDatasetType()
}

type Node interface {
	IsNode()
}

type AreaQuery struct {
	ParentCode   *AreaCode `json:"parentCode"`
	DatasetTypes []string  `json:"datasetTypes"`
	SearchTokens []string  `json:"searchTokens"`
}

type City struct {
	ID             ID          `json:"id"`
	Code           AreaCode    `json:"code"`
	Name           string      `json:"name"`
	PrefectureID   ID          `json:"prefecture_id"`
	PrefectureCode AreaCode    `json:"prefecture_code"`
	Prefecture     *Prefecture `json:"prefecture"`
	Wards          []*Ward     `json:"wards"`
	Datasets       []Dataset   `json:"datasets"`
}

func (City) IsArea() {}
func (City) IsNode() {}

type DatasetForAreaQuery struct {
	ExcludeTypes []string `json:"excludeTypes"`
	IncludeTypes []string `json:"includeTypes"`
	SearchTokens []string `json:"searchTokens"`
	// この地域の配下にある全ての自治体のデータセットを含めるかどうか。
	// 例えば、札幌市の場合、札幌市自体のデータだけでなく、札幌市の配下の全ての区（例えば中央区や北区）のデータセットも含めるかどうか。
	Deep *bool `json:"deep"`
}

type DatasetQuery struct {
	AreaCodes    []AreaCode `json:"areaCodes"`
	ExcludeTypes []string   `json:"excludeTypes"`
	IncludeTypes []string   `json:"includeTypes"`
	SearchTokens []string   `json:"searchTokens"`
	// areaCodesで指定された地域の配下にある全ての自治体のデータセットを含めるかどうか。
	// 例えば、札幌市を指定した場合、札幌市自体のデータだけでなく、札幌市の配下の全ての区（例えば中央区や北区）のデータセットも含めるかどうか。
	Deep *bool `json:"deep"`
}

type DatasetTypeQuery struct {
	Category    *DatasetTypeCategory `json:"category"`
	PlateauSpec *string              `json:"plateauSpec"`
	Year        *int                 `json:"year"`
}

type GenericDataset struct {
	ID             ID                    `json:"id"`
	Name           string                `json:"name"`
	Subname        *string               `json:"subname"`
	Description    *string               `json:"description"`
	Year           int                   `json:"year"`
	Groups         []string              `json:"groups"`
	PrefectureID   ID                    `json:"prefecture_id"`
	PrefectureCode AreaCode              `json:"prefecture_code"`
	CityID         *ID                   `json:"city_id"`
	CityCode       *AreaCode             `json:"city_code"`
	WardID         *ID                   `json:"ward_id"`
	WardCode       *AreaCode             `json:"ward_code"`
	Prefecture     *Prefecture           `json:"prefecture"`
	City           *City                 `json:"city"`
	Ward           *Ward                 `json:"ward"`
	TypeID         ID                    `json:"type_id"`
	Type           *GenericDatasetType   `json:"type"`
	Data           []*GenericDatasetItem `json:"data"`
}

func (GenericDataset) IsDataset() {}
func (GenericDataset) IsNode()    {}

type GenericDatasetItem struct {
	ID       ID              `json:"id"`
	Format   DatasetFormat   `json:"format"`
	Name     string          `json:"name"`
	URL      string          `json:"url"`
	Layers   []string        `json:"layers"`
	ParentID ID              `json:"parent_id"`
	Parent   *GenericDataset `json:"parent"`
}

func (GenericDatasetItem) IsDatasetItem() {}
func (GenericDatasetItem) IsNode()        {}

type GenericDatasetType struct {
	ID          ID                  `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	EnglishName string              `json:"englishName"`
	Category    DatasetTypeCategory `json:"category"`
}

func (GenericDatasetType) IsDatasetType() {}
func (GenericDatasetType) IsNode()        {}

// PLATEAU都市モデルの通常のデータセット。例えば、地物型が建築物モデル（bldg）などのデータセットです。
type PlateauDataset struct {
	ID             ID                    `json:"id"`
	Name           string                `json:"name"`
	Subname        *string               `json:"subname"`
	Description    *string               `json:"description"`
	Year           int                   `json:"year"`
	Groups         []string              `json:"groups"`
	PrefectureID   ID                    `json:"prefecture_id"`
	PrefectureCode AreaCode              `json:"prefecture_code"`
	CityID         *ID                   `json:"city_id"`
	CityCode       *AreaCode             `json:"city_code"`
	WardID         *ID                   `json:"ward_id"`
	WardCode       *AreaCode             `json:"ward_code"`
	Prefecture     *Prefecture           `json:"prefecture"`
	City           *City                 `json:"city"`
	Ward           *Ward                 `json:"ward"`
	TypeID         ID                    `json:"type_id"`
	Type           *PlateauDatasetType   `json:"type"`
	Data           []*PlateauDatasetItem `json:"data"`
}

func (PlateauDataset) IsDataset() {}
func (PlateauDataset) IsNode()    {}

type PlateauDatasetItem struct {
	ID       ID              `json:"id"`
	Format   DatasetFormat   `json:"format"`
	Name     string          `json:"name"`
	URL      string          `json:"url"`
	Layers   []string        `json:"layers"`
	ParentID ID              `json:"parent_id"`
	Parent   *PlateauDataset `json:"parent"`
	Lod      *float64        `json:"lod"`
	Texture  *Texture        `json:"texture"`
}

func (PlateauDatasetItem) IsDatasetItem() {}
func (PlateauDatasetItem) IsNode()        {}

type PlateauDatasetType struct {
	ID            ID                  `json:"id"`
	Code          string              `json:"code"`
	Name          string              `json:"name"`
	EnglishName   string              `json:"englishName"`
	Category      DatasetTypeCategory `json:"category"`
	PlateauSpec   *PlateauSpec        `json:"plateauSpec"`
	PlateauSpecID ID                  `json:"plateauSpecId"`
	Year          int                 `json:"year"`
	Flood         bool                `json:"flood"`
}

func (PlateauDatasetType) IsDatasetType() {}
func (PlateauDatasetType) IsNode()        {}

// PLATEAU都市モデルのデータセットのうち、地物型が洪水・高潮・津波・内水浸水想定区域モデル（fld, htd, tnm, ifld）のデータセットです。
type PlateauFloodingDataset struct {
	ID             ID                            `json:"id"`
	Name           string                        `json:"name"`
	Subname        *string                       `json:"subname"`
	Description    *string                       `json:"description"`
	Year           int                           `json:"year"`
	Groups         []string                      `json:"groups"`
	PrefectureID   ID                            `json:"prefecture_id"`
	PrefectureCode AreaCode                      `json:"prefecture_code"`
	CityID         *ID                           `json:"city_id"`
	CityCode       *AreaCode                     `json:"city_code"`
	WardID         *ID                           `json:"ward_id"`
	WardCode       *AreaCode                     `json:"ward_code"`
	Prefecture     *Prefecture                   `json:"prefecture"`
	City           *City                         `json:"city"`
	Ward           *Ward                         `json:"ward"`
	TypeID         ID                            `json:"type_id"`
	Type           *PlateauDatasetType           `json:"type"`
	Data           []*PlateauFloodingDatasetItem `json:"data"`
	// 河川。地物型が洪水浸水想定区域（fld）の場合のみ存在します。
	River *River `json:"river"`
}

func (PlateauFloodingDataset) IsDataset() {}
func (PlateauFloodingDataset) IsNode()    {}

type PlateauFloodingDatasetItem struct {
	ID       ID              `json:"id"`
	Format   DatasetFormat   `json:"format"`
	Name     string          `json:"name"`
	URL      string          `json:"url"`
	Layers   []string        `json:"layers"`
	ParentID ID              `json:"parent_id"`
	Parent   *PlateauDataset `json:"parent"`
	// 浸水規模
	FloodingScale FloodingScale `json:"floodingScale"`
}

func (PlateauFloodingDatasetItem) IsDatasetItem() {}
func (PlateauFloodingDatasetItem) IsNode()        {}

type PlateauSpec struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

func (PlateauSpec) IsNode() {}

type Prefecture struct {
	ID       ID        `json:"id"`
	Code     AreaCode  `json:"code"`
	Name     string    `json:"name"`
	Cities   []*City   `json:"cities"`
	Datasets []Dataset `json:"datasets"`
}

func (Prefecture) IsArea() {}
func (Prefecture) IsNode() {}

type RelatedDataset struct {
	ID             ID                    `json:"id"`
	Name           string                `json:"name"`
	Subname        *string               `json:"subname"`
	Description    *string               `json:"description"`
	Year           int                   `json:"year"`
	Groups         []string              `json:"groups"`
	PrefectureID   ID                    `json:"prefecture_id"`
	PrefectureCode AreaCode              `json:"prefecture_code"`
	CityID         *ID                   `json:"city_id"`
	CityCode       *AreaCode             `json:"city_code"`
	WardID         *ID                   `json:"ward_id"`
	WardCode       *AreaCode             `json:"ward_code"`
	Prefecture     *Prefecture           `json:"prefecture"`
	City           *City                 `json:"city"`
	Ward           *Ward                 `json:"ward"`
	TypeID         ID                    `json:"type_id"`
	Type           *RelatedDatasetType   `json:"type"`
	Data           []*RelatedDatasetItem `json:"data"`
}

func (RelatedDataset) IsDataset() {}
func (RelatedDataset) IsNode()    {}

type RelatedDatasetItem struct {
	ID       ID              `json:"id"`
	Format   DatasetFormat   `json:"format"`
	Name     string          `json:"name"`
	URL      string          `json:"url"`
	Layers   []string        `json:"layers"`
	ParentID ID              `json:"parent_id"`
	Parent   *RelatedDataset `json:"parent"`
}

func (RelatedDatasetItem) IsDatasetItem() {}
func (RelatedDatasetItem) IsNode()        {}

type RelatedDatasetType struct {
	ID          ID                  `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	EnglishName string              `json:"englishName"`
	Category    DatasetTypeCategory `json:"category"`
}

func (RelatedDatasetType) IsDatasetType() {}
func (RelatedDatasetType) IsNode()        {}

// 河川
type River struct {
	// 河川名。通常「〜水系〜川」という形式になります。
	Name string `json:"name"`
	// 管理区間
	Admin RiverAdmin `json:"admin"`
}

type Ward struct {
	ID             ID          `json:"id"`
	Code           AreaCode    `json:"code"`
	Name           string      `json:"name"`
	PrefectureID   ID          `json:"prefecture_id"`
	PrefectureCode AreaCode    `json:"prefecture_code"`
	CityID         ID          `json:"city_id"`
	CityCode       AreaCode    `json:"city_code"`
	Prefecture     *Prefecture `json:"prefecture"`
	City           *City       `json:"city"`
	Datasets       []Dataset   `json:"datasets"`
}

func (Ward) IsArea() {}
func (Ward) IsNode() {}

type DatasetFormat string

const (
	DatasetFormatCSV           DatasetFormat = "CSV"
	DatasetFormatCzml          DatasetFormat = "CZML"
	DatasetFormatCesium3DTiles DatasetFormat = "Cesium3DTiles"
	DatasetFormatGltf          DatasetFormat = "GLTF"
	DatasetFormatGTFSRelatime  DatasetFormat = "GTFSRelatime"
	DatasetFormatGeoJSON       DatasetFormat = "GeoJSON"
	DatasetFormatMvt           DatasetFormat = "MVT"
	DatasetFormatTms           DatasetFormat = "TMS"
	DatasetFormatTiles         DatasetFormat = "Tiles"
	DatasetFormatWms           DatasetFormat = "WMS"
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

type DatasetTypeCategory string

const (
	DatasetTypeCategoryPlateau DatasetTypeCategory = "Plateau"
	DatasetTypeCategoryRelated DatasetTypeCategory = "Related"
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

// 浸水規模
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

type Texture string

const (
	TextureNone           Texture = "NONE"
	TextureLowResolution  Texture = "LOW_RESOLUTION"
	TextureHighResolution Texture = "HIGH_RESOLUTION"
)

var AllTexture = []Texture{
	TextureNone,
	TextureLowResolution,
	TextureHighResolution,
}

func (e Texture) IsValid() bool {
	switch e {
	case TextureNone, TextureLowResolution, TextureHighResolution:
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
