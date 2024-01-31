package geospatialjpv3

import (
	"context"
	"fmt"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func (h *handler) Publish(ctx context.Context, w *cmswebhook.Payload) error {
	cms := h.cms
	log.Infofc(ctx, "geospatialjpv3: publish: %+v", w)

	cityItemRaw := w.ItemData.Item
	cityItem := CityItemFrom(cityItemRaw)
	pkgYear := cityItem.YearInt()

	if cityItem.CityCode == "" || cityItem.CityName == "" || cityItem.CityNameEn == "" || pkgYear == 0 {
		return fmt.Errorf("invalid city item")
	}

	log.Debugfc(ctx, "geospatialjpv3: cityItem: %s", ppp.Sprint(cityItem))

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

	log.Debugfc(ctx, "geospatialjpv3: pkg: %s", ppp.Sprint(pkg))
	seed, err := getSeed(ctx, cms, cityItem)
	if err != nil {
		return fmt.Errorf("failed to get seed: %w", err)
	}

	if seed.CityGML != "" {
		log.Debugfc(ctx, "geospatialjpv3: citygml: %s", seed.CityGML)
		if err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf("CityGML（v%d）", seed.Version),
			URL:         seed.CityGML,
			Description: "",
		}); err != nil {
			return fmt.Errorf("failed to create or update resource (citygml): %w", err)
		}
	}

	if seed.Plateau != "" {
		log.Debugfc(ctx, "geospatialjpv3: plateau: %s", seed.Plateau)
		if err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf("3D Tiles, MVT（%d）", seed.Version),
			URL:         seed.Plateau,
			Description: "",
		}); err != nil {
			return fmt.Errorf("failed to create or update resource (plateau): %w", err)
		}
	}

	if seed.Related != "" {
		log.Debugfc(ctx, "geospatialjpv3: related: %s", seed.Related)
		if err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf(("関連データセット（%d）"), seed.Version),
			URL:         seed.Related,
			Description: "",
		}); err != nil {
			return fmt.Errorf("failed to create or update resource (related): %w", err)
		}
	}

	return nil
}
