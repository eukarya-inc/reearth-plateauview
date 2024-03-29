package cmsintegrationv2

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

type Status string

const (
	StatusReady      Status = "未実行"
	StatusProcessing Status = "実行中"
	StatusOK         Status = "完了"
	StatusError      Status = "エラー"
)

type Conversion string

const (
	ConversionDisabled Conversion = "変換しない"
	ConversionEnabled  Conversion = "変換する"
)

func (s Conversion) Enabled() bool {
	return s == ConversionEnabled
}

type Separation string

func (s Separation) Enabled() bool {
	return string(s) != "分割しない"
}

type PRCS = cmsintegrationcommon.PRCS

type Item struct {
	ID string `json:"id,omitempty" cms:"id"`
	// select: specification
	Specification string `json:"specification,omitempty" cms:"specification,select"`
	// asset: citygml
	CityGML string `json:"citygml,omitempty" cms:"citygml,asset"`
	// select: conversion_enabled: 変換する, 変換しない
	ConversionEnabled Conversion `json:"conversion_enabled,omitempty" cms:"conversion_enabled,select"`
	// select: prcs: 第1系~第19系
	PRCS PRCS `json:"prcs" cms:"prcs,select"`
	// asset: quality_check_params
	QualityCheckParams string `json:"quality_check_params,omitempty" cms:"quality_check_params,asset"`
	// select: devide_odc: 分割する, 分割しない
	DevideODC Separation `json:"devide_odc,omitempty" cms:"devide_odc,select"`
	// asset[]: bldg
	Bldg []string `json:"bldg,omitempty" cms:"bldg,asset"`
	// asset: tran
	Tran []string `json:"tran,omitempty" cms:"tran,asset"`
	// asset: frn
	Frn []string `json:"frn,omitempty" cms:"frn,asset"`
	// asset: veg
	Veg []string `json:"veg,omitempty" cms:"veg,asset"`
	// asset: luse
	Luse []string `json:"luse,omitempty" cms:"luse,asset"`
	// asset: lsld
	Lsld []string `json:"lsld,omitempty" cms:"lsld,asset"`
	// asset: urf
	Urf []string `json:"urf,omitempty" cms:"urf,asset"`
	// asset[]: fld
	Fld []string `json:"fld,omitempty" cms:"fld,asset"`
	// asset[]: tnm
	Tnm []string `json:"tnm,omitempty" cms:"tnm,asset"`
	// asset[]: htd
	Htd []string `json:"htd,omitempty" cms:"htd,asset"`
	// asset[]: ifld
	Ifld []string `json:"ifld,omitempty" cms:"ifld,asset"`
	// asset: all
	All string `json:"all,omitempty" cms:"all,asset"`
	// asset: dictionary
	Dictionary string `json:"dictionary,omitempty" cms:"dictionary,asset"`
	// textarea: dic
	Dic string `json:"dic,omitempty" cms:"dic,textarea"`
	// select: conversion_status: 未実行, 実行中, 完了, エラー
	ConversionStatus Status `json:"conversion_status,omitempty" cms:"conversion_status,select"`

	// SDK
	// asset: max_lod
	MaxLOD string `json:"max_lod,omitempty" cms:"max_lod,asset"`
	// select: max_lod_status: 未実行, 実行中, 完了, エラー
	MaxLODStatus Status `json:"max_lod_status,omitempty" cms:"max_lod_status,select"`
	// select: sdk_publication: 公開する・公開しない
	SDKPublication string `json:"sdk_publication,omitempty" cms:"sdk_publication,select"`
	// select: dem: 無し・有り
	Dem       string `json:"dem,omitempty" cms:"dem,select"`
	ProjectID string `json:"-" cms:"-"`
}

func (i Item) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}

func ItemFrom(item cms.Item) (i Item) {
	item.Unmarshal(&i)
	return
}

func (i Item) IsPublicOnSDK() bool {
	return i.SDKPublication == "公開する"
}

func (i Item) HasDem() bool {
	return i.Dem == "有り"
}
