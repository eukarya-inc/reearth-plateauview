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
	Tnm    *AssetNameExTnm
	Ex     string
}

func (n AssetNameEx) Type() string {
	switch {
	case n.Normal != nil:
		return n.Normal.Type
	case n.Urf != nil:
		return n.Urf.Type
	case n.Fld != nil:
		return n.Fld.Type
	case n.Tnm != nil:
		return n.Tnm.Type
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

type AssetNameExUrf struct {
	Type   string
	Name   string
	Format string
	LOD    int
}

type AssetNameExFld struct {
	Type      string
	Admin     string
	River     string
	Format    string
	L         int
	NoTexture bool
}

type AssetNameExTnm struct {
	Type      string
	Name      string
	Format    string
	NoTexture bool
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

	ex.Tnm = ParseAssetNameExTnm(name)
	if ex.Tnm != nil {
		return
	}

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

var reAasetNameExUrf = regexp.MustCompile(`^([a-z]+)_([A-Za-z0-9-_]+)_(mvt|3dtiles)(_lod\d+)?$`)

func ParseAssetNameExUrf(name string) *AssetNameExUrf {
	if name == "" {
		return nil
	}

	m := reAasetNameExUrf.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	lod := 0
	if m[4] != "" {
		lod, _ = strconv.Atoi(m[4][4:])
	}

	return &AssetNameExUrf{
		Type:   m[1],
		Name:   m[2],
		Format: m[3],
		LOD:    lod,
	}
}

var reAasetNameExFld = regexp.MustCompile(`^fld_([a-z0-9-]+)_([a-z0-9-_]+)_3dtiles_(l\d+)(_no_texture)?$`)

func ParseAssetNameExFld(name string) *AssetNameExFld {
	if name == "" {
		return nil
	}

	m := reAasetNameExFld.FindStringSubmatch(name)
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

var reAasetNameExTnm = regexp.MustCompile(`^(tnm|htd|ifld)_([a-z0-9-_]+)_3dtiles(_no_texture)?$`)

func ParseAssetNameExTnm(name string) *AssetNameExTnm {
	if name == "" {
		return nil
	}

	m := reAasetNameExTnm.FindStringSubmatch(name)
	if len(m) == 0 {
		return nil
	}

	return &AssetNameExTnm{
		Type:      m[1],
		Name:      m[2],
		Format:    "3dtiles",
		NoTexture: m[3] != "",
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
