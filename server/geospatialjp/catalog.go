package geospatialjp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type Catalog struct {
	// タイトル
	Title string `json:"title,omitempty"`
	// URL
	URL string `json:"url,omitempty"`
	// 説明
	Description string `json:"description,omitempty"`
	// タグ
	Tags []string `json:"tags,omitempty"`
	// ライセンス
	License string `json:"license,omitempty"`
	// 組織
	Organization string `json:"organization,omitempty"`
	// 公開・非公開
	Public string `json:"public,omitempty"`
	// ソース
	Source string `json:"source,omitempty"`
	// バージョン
	Version string `json:"version,omitempty"`
	// 作成者
	Author string `json:"author,omitempty"`
	// 作成者のメールアドレス
	AuthorEmail string `json:"authorEmail,omitempty"`
	// メンテナー（保守者）
	Maintainer string `json:"maintainer,omitempty"`
	// メンテナー（保守者）のメールアドレス
	MaintainerEmail string `json:"maintainerEmail,omitempty"`
	// spatial*
	Spatial string `json:"spatial,omitempty"`
	// データ品質
	DataQuality string `json:"dataQuality,omitempty"`
	// 制約
	Constraints string `json:"constraints,omitempty"`
	// データ登録日
	RegisteredAt string `json:"registeredAt,omitempty"`
	// 有償無償区分*
	FreeOrProvidedClassification string `json:"freeOrProvidedClassification,omitempty"`
	// 災害時区分*
	DisasterClassification string `json:"disasterClassification,omitempty"`
	// 地理的範囲
	GeoArea string `json:"geoArea,omitempty"`
	// サムネイル画像
	Thumbnail []byte `json:"-"`
	// サムネイル画像のファイル名
	ThumbnailFileName string `json:"-"`
	// 価格情報
	Price string `json:"price,omitempty"`
	// 使用許諾
	LicenseAgreement string `json:"licenseAgreement,omitempty"`
	// カスタムフィールド
	CustomFields map[string]any `json:"customFields,omitempty"`
}

type CatalogFile struct {
	file *excelize.File
}

func NewCatalogFile(file *excelize.File) *CatalogFile {
	return &CatalogFile{
		file: file,
	}
}

func (c *CatalogFile) Parse() (res Catalog, err error) {
	sheet, err := c.getSheet()
	if err != nil {
		return res, err
	}

	errs := []error{}

	res.Title, errs = c.getCellValue(sheet, "タイトル", "D2", errs)
	res.URL, errs = c.getCellValue(sheet, "URL", "D3", errs)
	res.Description, errs = c.getCellValue(sheet, "説明", "D4", errs)
	res.Tags, errs = c.getCellValueAsTags(sheet, "タグ", "D5", errs)
	res.License, errs = c.getCellValue(sheet, "ライセンス", "D6", errs)
	res.Organization, errs = c.getCellValue(sheet, "組織", "D7", errs)
	res.Public, errs = c.getCellValue(sheet, "公開・非公開", "D8", errs)
	res.Source, errs = c.getCellValue(sheet, "ソース", "D9", errs)
	res.Version, errs = c.getCellValue(sheet, "バージョン", "D10", errs)
	res.Author, errs = c.getCellValue(sheet, "作成者", "D11", errs)
	res.AuthorEmail, errs = c.getCellValue(sheet, "作成者のメールアドレス", "D12", errs)
	res.Maintainer, errs = c.getCellValue(sheet, "メンテナー（保守者）", "D13", errs)
	res.MaintainerEmail, errs = c.getCellValue(sheet, "メンテナー（保守者）のメールアドレス", "D14", errs)
	res.Spatial, errs = c.getCellValue(sheet, "spatial*", "D15", errs)
	res.DataQuality, errs = c.getCellValue(sheet, "データ品質", "D16", errs)
	res.Constraints, errs = c.getCellValue(sheet, "制約", "D17", errs)
	res.RegisteredAt, errs = c.getCellValue(sheet, "データ登録日", "D18", errs)
	res.FreeOrProvidedClassification, errs = c.getCellValue(sheet, "有償無償区分*", "D19", errs)
	res.DisasterClassification, errs = c.getCellValue(sheet, "災害時区分*", "D20", errs)
	res.GeoArea, errs = c.getCellValue(sheet, "地理的範囲", "D21", errs)
	res.ThumbnailFileName, res.Thumbnail, errs = c.getPicture(sheet, "サムネイル画像", "D22", errs)
	res.Price, errs = c.getCellValue(sheet, "価格情報", "D23", errs)
	res.LicenseAgreement, errs = c.getCellValue(sheet, "使用許諾", "D24", errs)
	// TODO: メタデータ is not implemented yet

	if len(errs) > 0 {
		return res, fmt.Errorf("目録の読み込みに失敗しました。%w", errorsJoin(errs))
	}
	return res, nil
}

func (c *CatalogFile) DeleteSheet() error {
	sheet, err := c.getSheet()
	if err != nil {
		return err
	}
	c.file.DeleteSheet(sheet)
	return nil
}

func (c *CatalogFile) File() *excelize.File {
	return c.file
}

func (c *CatalogFile) getSheet() (string, error) {
	if i := c.file.GetSheetIndex("G空間登録用メタデータ "); i < 0 {
		if i = c.file.GetSheetIndex("G空間登録用メタデータ"); i < 0 {
			return "", errors.New("シート「G空間登録用メタデータ」が見つかりませんでした。")
		}
		return "G空間登録用メタデータ", nil
	}
	return "G空間登録用メタデータ ", nil
}

func (c *CatalogFile) getCellValue(sheet, name, _axis string, errs []error) (string, []error) {
	pos, errs := c.findCell(sheet, name, errs)
	if pos != "" {
		cp, err := ParseCellPos(pos)
		if err != nil {
			return "", append(errs, err)
		}
		pos = cp.ShiftX(2).String()
	}

	cell, err := c.file.GetCellValue(sheet, pos)
	if err != nil {
		return "", append(errs, fmt.Errorf("「%s」が見つかりませんでした。", name))
	}
	return strings.ReplaceAll(cell, "\u2028", "\n"), nil
}

func (c *CatalogFile) getCellValueAsTags(sheet, name, axis string, errs []error) ([]string, []error) {
	cell, errs := c.getCellValue(sheet, name, axis, errs)
	tags := lo.Map(
		lo.FlatMap(
			strings.Split(cell, ","),
			func(s string, _ int) []string {
				return strings.Split(s, "、")
			}),
		func(s string, _ int) string {
			return strings.TrimSpace(s)
		},
	)
	return tags, errs
}

func (c *CatalogFile) getPicture(sheet, name, axis string, errs []error) (string, []byte, []error) {
	file, raw, err := c.file.GetPicture(sheet, axis)
	if err != nil {
		return "", nil, append(errs, fmt.Errorf("「%s」が見つかりませんでした。", name))
	}
	return file, raw, errs
}

func (c *CatalogFile) findCell(sheet, name string, errs []error) (string, []error) {
	pos, err := c.file.SearchSheet(sheet, name)
	if err != nil || len(pos) == 0 {
		return "", append(errs, fmt.Errorf("「%s」が見つかりませんでした。", name))
	}
	return minXPos(pos), errs
}

func errorsJoin(errs []error) error {
	return errors.New(strings.Join(lo.Map(errs, func(e error, _ int) string { return e.Error() }), ""))
}

func minXPos(pos []string) string {
	pos2, err := util.TryMap(pos, ParseCellPos)
	if err != nil {
		return ""
	}

	minPos := lo.MinBy(pos2, func(a, b CellPos) bool {
		return a.x < b.x
	})
	return minPos.String()
}
