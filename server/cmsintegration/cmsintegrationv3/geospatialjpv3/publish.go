package geospatialjpv3

import (
	"context"
	"fmt"

	"github.com/k0kubun/pp/v3"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func (h *handler) Publish(ctx context.Context, w *cmswebhook.Payload) error {
	const pkgYear = 2023
	cms := h.cms
	log.Infofc(ctx, "geospatialjpv3: publish: %+v", w)

	cityItemRaw, err := cms.GetItem(ctx, w.ItemData.Item.ID, false)
	if err != nil {
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	{
		pp := pp.New()
		pp.SetColoringEnabled(false)
		s := pp.Sprint(cityItem)
		log.Debugfc(ctx, "geospatialjpv3: cityItem: %s", s)
	}

	pkg, err := h.findOrCreatePackage(
		ctx,
		cityItem.CityCode,
		cityItem.CityNameEn,
		pkgYear,
		CityInfo{NameJa: cityItem.CityName},
	)
	if err != nil {
		return fmt.Errorf("failed to find or create package: %w", err)
	}
	log.Debugfc(ctx, "geospatialjpv3: pkg: %+v", pkg)

	geoItem, err := getGeospatialItems(ctx, cms, cityItem)
	if err != nil {
		return fmt.Errorf("failed to get all feature items: %w", err)
	}

	if geoItem.CityGML != "" {
		log.Debugfc(ctx, "geospatialjpv3: citygml: %s", geoItem.CityGML)
		h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        "CityGML（v3）",
			URL:         geoItem.CityGML,
			Description: "",
		})
	}

	if geoItem.Plateau != "" {
		log.Debugfc(ctx, "geospatialjpv3: plateau: %s", geoItem.Plateau)
		h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        "3D Tiles, MVT（v3）",
			URL:         geoItem.Plateau,
			Description: "",
		})
	}

	// if geoItem.Related != "" {
	// 	log.Debugfc(ctx, "geospatialjpv3: related: %s", geoItem.Related)
	// 	h.createOrUpdateResource(ctx, pkg, ResourceInfo{
	// 		Name:        "関連データ",
	// 		URL:         geoItem.Related,
	// 		Description: "",
	// 	})
	// }

	return nil
}
