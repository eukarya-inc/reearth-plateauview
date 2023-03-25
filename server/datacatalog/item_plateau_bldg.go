package datacatalog

import (
	"fmt"
	"strconv"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/samber/lo"
)

func (i PlateauItem) BldgItems(c PlateauIntermediateItem) []*DataCatalogItem {
	assets := assetsByWards(i.Bldg)
	if len(assets) == 0 {
		return nil
	}

	firstCode := lo.Min(lo.Filter(lo.MapToSlice(assets, func(k string, v []*cms.PublicAsset) int {
		if len(v) == 0 {
			return 0
		}
		an := AssetNameFrom(v[0].URL)
		wc, _ := strconv.Atoi(an.WardCode)
		return wc
	}), func(i int, _ int) bool { return i > 0 }))

	return lo.Filter(lo.MapToSlice(assets, func(k string, v []*cms.PublicAsset) *DataCatalogItem {
		s := BldgSetFrom(v)
		if s == nil || s.MaxLOD.Texture == nil {
			return nil
		}

		an := AssetNameFrom(s.MaxLOD.Texture.URL)
		wc, _ := strconv.Atoi(an.WardCode)
		dci := c.DataCatalogItem(
			"建築物モデル",
			an,
			s.MaxLOD.Texture.URL,
			i.DescriptionBldg,
			nil,
			firstCode > 0 && firstCode == wc,
		)

		if s.MaxLOD.LowTexture != nil {
			dci.BldgLowTextureURL = assetURLFromFormat(s.MaxLOD.LowTexture.URL, "3dtiles")
		}

		if s.MaxLOD.NoTexture != nil {
			dci.BldgNoTextureURL = assetURLFromFormat(s.MaxLOD.NoTexture.URL, "3dtiles")
		}

		dci.SearchIndex = searchIndexURLFrom(i.SearchIndex, dci.WardCode)
		dci.Config = s.Config()

		return dci
	}), func(a *DataCatalogItem, _ int) bool {
		return a != nil
	})
}

type BldgSet struct {
	MaxLOD *BldgSetLOD
	LOD0   *BldgSetLOD
	LOD1   *BldgSetLOD
	LOD2   *BldgSetLOD
	LOD3   *BldgSetLOD
	LOD4   *BldgSetLOD
}

type BldgSetLOD struct {
	LOD        int
	Texture    *cms.PublicAsset
	LowTexture *cms.PublicAsset
	NoTexture  *cms.PublicAsset
}

func BldgSetFrom(a []*cms.PublicAsset) *BldgSet {
	lods, maxlod := assetWithLODFromList(a)
	if len(lods) == 0 {
		return nil
	}
	return &BldgSet{
		MaxLOD: bldgSetLODFrom(lods, maxlod),
		LOD0:   bldgSetLODFrom(lods, 0),
		LOD1:   bldgSetLODFrom(lods, 1),
		LOD2:   bldgSetLODFrom(lods, 2),
		LOD3:   bldgSetLODFrom(lods, 3),
		LOD4:   bldgSetLODFrom(lods, 4),
	}
}

func bldgSetLODFrom(assets []assetWithLOD, lod int) *BldgSetLOD {
	tex, _ := lo.Find(assets, func(a assetWithLOD) bool {
		return a.LOD == lod && !a.F.LowTexture && !a.F.NoTexture
	})
	lowtex, _ := lo.Find(assets, func(a assetWithLOD) bool {
		return a.LOD == lod && a.F.LowTexture
	})
	notex, _ := lo.Find(assets, func(a assetWithLOD) bool {
		return a.LOD == lod && a.F.NoTexture
	})

	if tex.A == nil && lowtex.A == nil && notex.A == nil {
		return nil
	}

	return &BldgSetLOD{
		LOD:        lod,
		Texture:    tex.A,
		LowTexture: lowtex.A,
		NoTexture:  notex.A,
	}
}

func (s BldgSet) Config() (c DataCatalogItemConfig) {
	if l := s.LOD0.Config(); len(l) > 0 {
		c.Data = append(c.Data, l...)
	}
	if l := s.LOD1.Config(); len(l) > 0 {
		c.Data = append(c.Data, l...)
	}
	if l := s.LOD2.Config(); len(l) > 0 {
		c.Data = append(c.Data, l...)
	}
	if l := s.LOD3.Config(); len(l) > 0 {
		c.Data = append(c.Data, l...)
	}
	if l := s.LOD4.Config(); len(l) > 0 {
		c.Data = append(c.Data, l...)
	}
	return
}

func (s *BldgSetLOD) Config() (c []DataCatalogItemConfigItem) {
	if s == nil {
		return nil
	}

	if s.Texture != nil {
		c = append(c, DataCatalogItemConfigItem{
			Name: fmt.Sprintf("LOD%d", s.LOD),
			URL:  assetURLFromFormat(s.Texture.URL, "3dtiles"),
			Type: "3dtiles",
		})
	}

	if s.LowTexture != nil {
		c = append(c, DataCatalogItemConfigItem{
			Name: fmt.Sprintf("LOD%d（低解像度）", s.LOD),
			URL:  assetURLFromFormat(s.LowTexture.URL, "3dtiles"),
			Type: "3dtiles",
		})
	}

	if s.NoTexture != nil {
		c = append(c, DataCatalogItemConfigItem{
			Name: fmt.Sprintf("LOD%d（テクスチャなし）", s.LOD),
			URL:  assetURLFromFormat(s.NoTexture.URL, "3dtiles"),
			Type: "3dtiles",
		})
	}

	return
}
