package geospatialjpv3

import (
	"context"
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
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

	resources := []ckan.Resource{}

	log.Debugfc(ctx, "geospatialjpv3: seed: %s", ppp.Sprint(seed))
	if seed.CityGML != "" {
		log.Debugfc(ctx, "geospatialjpv3: citygml: %s", seed.CityGML)
		r, err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf("CityGML（v%d）", seed.Version),
			URL:         seed.CityGML,
			Description: "",
		})
		if err != nil {
			return fmt.Errorf("failed to create or update resource (citygml): %w", err)
		}
		resources = append(resources, r)
	}

	if seed.Plateau != "" {
		log.Debugfc(ctx, "geospatialjpv3: plateau: %s", seed.Plateau)
		r, err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf("3D Tiles, MVT（v%d）", seed.Version),
			URL:         seed.Plateau,
			Description: "",
		})
		if err != nil {
			return fmt.Errorf("failed to create or update resource (plateau): %w", err)
		}
		resources = append(resources, r)
	}

	if seed.Related != "" {
		log.Debugfc(ctx, "geospatialjpv3: related: %s", seed.Related)
		r, err := h.createOrUpdateResource(ctx, pkg, ResourceInfo{
			Name:        fmt.Sprintf(("関連データセット（v%d）"), seed.Version),
			URL:         seed.Related,
			Description: "",
		})
		if err != nil {
			return fmt.Errorf("failed to create or update resource (related): %w", err)
		}
		resources = append(resources, r)
	}

	if (seed.CityGML != "" || seed.Plateau != "" || seed.Related != "") && shouldReorder(pkg, seed.Version) {
		log.Debugfc(ctx, "geospatialjpv3: reorder: %v", resources)
		resourceIDs := lo.Map(resources, func(r ckan.Resource, _ int) string {
			return r.ID
		})

		if err := h.reorderResources(ctx, pkg.ID, resourceIDs); err != nil {
			return fmt.Errorf("failed to reorder resources: %w", err)
		}
	}

	return nil
}

func shouldReorder(pkg *ckan.Package, currentVersion int) bool {
	for _, res := range pkg.Resources {
		if strings.Contains(res.Name, fmt.Sprintf("(v%d)", currentVersion)) ||
			strings.Contains(res.Name, fmt.Sprintf("（v%d）", currentVersion)) {
			return false
		}
	}
	return true
}
