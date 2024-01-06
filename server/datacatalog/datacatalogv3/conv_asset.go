package datacatalogv3

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/samber/lo"
)

type AssetName struct {
	CityCode    string
	CityName    string
	Provider    string
	Year        int
	Format      string
	UpdateCount int
	Ex          AssetNameEx
}

func (n AssetName) String() string {
	var ex string
	if n.Ex.Ex != "" {
		ex = "_" + n.Ex.Ex
	}
	return fmt.Sprintf("%s_%s_%s_%d_%s_%d_op%s", n.CityCode, n.CityName, n.Provider, n.Year, n.Format, n.UpdateCount, ex)
}

type AssetNameEx struct {
	Normal *AssetNameExNormal
	Urf    *AssetNameExUrf
	Fld    *AssetNameExFld
	Ex     string
}

func (ex AssetNameEx) String() string {
	return ex.Ex
}

func (ex AssetNameEx) IsValid() bool {
	return ex.Normal != nil || ex.Urf != nil || ex.Fld != nil
}

func (ex AssetNameEx) Key() string {
	switch {
	case ex.Normal != nil:
		return ex.Normal.Key()
	case ex.Urf != nil:
		return ex.Urf.Key()
	case ex.Fld != nil:
		return ex.Fld.Key()
	}
	return ""
}

func (ex AssetNameEx) ItemKey() string {
	switch {
	case ex.Normal != nil:
		return ex.Normal.ItemKey()
	case ex.Urf != nil:
		return ex.Urf.ItemKey()
	case ex.Fld != nil:
		return ex.Fld.ItemKey()
	}
	return ""
}

func (ex AssetNameEx) DicKey() string {
	switch {
	case ex.Normal != nil:
		return ex.Normal.DicKey()
	case ex.Urf != nil:
		return ex.Urf.DicKey()
	case ex.Fld != nil:
		return ex.Fld.DicKey()
	}
	return ""
}

type AssetNameExNormal struct {
	Type      string
	Format    string
	WardCode  string
	WardName  string
	LOD       int
	NoTexture bool
}

func (ex AssetNameExNormal) Key() string {
	return ""
}

func (ex AssetNameExNormal) ItemKey() string {
	return ""
}

func (ex AssetNameExNormal) DicKey() string {
	return ""
}

type AssetNameExUrf struct {
	Type      string
	Name      string
	Format    string
	LOD       int
	NoTexture bool
}

func (ex AssetNameExUrf) Key() string {
	return ex.Name
}

func (ex AssetNameExUrf) ItemKey() string {
	return ex.Name
}

func (ex AssetNameExUrf) DicKey() string {
	return ex.Name
}

type AssetNameExFld struct {
	Type      string
	Admin     string
	River     string
	Format    string
	L         int
	NoTexture bool
}

func (ex AssetNameExFld) Key() string {
	return fmt.Sprintf("l%d", ex.L)
}

func (ex AssetNameExFld) ItemKey() string {
	return fmt.Sprintf("%s_%s", ex.Admin, ex.River)
}

func (ex AssetNameExFld) DicKey() string {
	return fmt.Sprintf("%s_l%d", ex.River, ex.L)
}

var reAssetName = regexp.MustCompile(`^(\d{5})_([a-z0-9-]+)_([a-z0-9-]+)_(\d{4})_(.+?)_(\d+)(?:_op$?)?(?:_(.+))?$`)

func ParseAssetName(name string) *AssetName {
	m := reAssetName.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	year, _ := strconv.Atoi(m[4])
	updateCount, _ := strconv.Atoi(m[6])
	var ex string
	if len(m) > 7 {
		ex = m[7]
	}

	return &AssetName{
		CityCode:    m[1],
		CityName:    m[2],
		Provider:    m[3],
		Year:        year,
		Format:      m[5],
		UpdateCount: updateCount,
		Ex:          ParseAssetNameEx(ex),
	}
}

func ParseAssetNameEx(name string) (ex AssetNameEx) {
	ex.Ex = name

	ex.Fld = ParseAssetNameExFld(name)
	if ex.Fld != nil {
		return
	}

	ex.Urf = ParseAssetNameExUrf(name)
	if ex.Urf != nil {
		return
	}

	ex.Normal = ParseAssetNameExNormal(name)
	return
}

var reAasetNameExNormal = regexp.MustCompile(`^([a-z]+)_(mvt|3dtiles)(?:_(\d+)_([a-z0-9-]+))?(_lod\d+)?(_no_texture)?$`)

func ParseAssetNameExNormal(name string) *AssetNameExNormal {
	if name == "" {
		return nil
	}

	m := reAasetNameExNormal.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	lod := 0
	if m[5] != "" {
		lod, _ = strconv.Atoi(m[5][4:])
	}

	return &AssetNameExNormal{
		Type:      m[1],
		Format:    m[2],
		WardCode:  m[3],
		WardName:  m[4],
		LOD:       lod,
		NoTexture: m[6] != "",
	}
}

var reAssetNameExUrf = regexp.MustCompile(`^([a-z]+)_([A-Za-z0-9-_]+)_(mvt|3dtiles)(_lod\d+)?(_no_texture)?$`)

func ParseAssetNameExUrf(name string) *AssetNameExUrf {
	if name == "" {
		return nil
	}

	m := reAssetNameExUrf.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	lod := 0
	if m[4] != "" {
		lod, _ = strconv.Atoi(m[4][4:])
	}

	return &AssetNameExUrf{
		Type:      m[1],
		Name:      m[2],
		Format:    m[3],
		LOD:       lod,
		NoTexture: m[5] != "",
	}
}

var reAssetNameExFld = regexp.MustCompile(`^fld_(natl|pref)_([a-z0-9-_]+)_3dtiles_(l\d+)(_no_texture)?$`)

func ParseAssetNameExFld(name string) *AssetNameExFld {
	if name == "" {
		return nil
	}

	m := reAssetNameExFld.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	l, _ := strconv.Atoi(m[3][1:])

	return &AssetNameExFld{
		Type:      "fld",
		Admin:     m[1],
		River:     m[2],
		Format:    "3dtiles",
		L:         l,
		NoTexture: m[4] != "",
	}
}

func ParseAssetUrls(urls []string) []*AssetName {
	return lo.Map(urls, func(u string, _ int) *AssetName {
		return ParseAssetName(nameWithoutExt(nameFromURL(u)))
	})
}

type RelatedAssetName struct {
	Code     string
	Name     string
	Year     int
	Provider string
	Type     string
	Format   string
}

var reRelatedAssetName = regexp.MustCompile(`^(\d{5})_([a-zA-Z0-9-]+)_([a-zA-Z0-9-]+)_(\d+)_([a-zA-Z0-9-]+)\.([a-z0-9]+)$`)

func ParseRelatedAssetName(name string) *RelatedAssetName {
	m := reRelatedAssetName.FindStringSubmatch(name)
	if m == nil {
		return nil
	}

	y, _ := strconv.Atoi(m[4])

	return &RelatedAssetName{
		Code:     m[1],
		Name:     m[2],
		Provider: m[3],
		Year:     y,
		Type:     m[5],
		Format:   m[6],
	}
}