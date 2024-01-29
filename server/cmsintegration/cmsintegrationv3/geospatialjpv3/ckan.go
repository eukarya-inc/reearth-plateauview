package geospatialjpv3

import (
	"context"
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

type CityInfo struct {
	NameJa string
}

func (s *handler) findOrCreatePackage(ctx context.Context, cityCode, cityNameEn string, year int, city CityInfo) (*ckan.Package, error) {
	// find
	pkg, pkgName, err := s.findPackage(ctx, cityCode, cityNameEn, year)
	if err != nil {
		return nil, fmt.Errorf("G空間情報センターからデータセットを検索できませんでした: %w", err)
	}

	// create
	if pkg == nil {
		newpkg := ckan.Package{
			Name:     pkgName,
			OwnerOrg: s.ckanOrg,
			Title:    fmt.Sprintf("3D都市モデル（Project PLATEAU）%s（%d年度）", city.NameJa, year),
			// TODO: use CityInfo
		}
		log.Infofc(ctx, "geospartialjp: package %s not found so new package will be created", pkgName)

		pkg2, err := s.ckan.CreatePackage(ctx, newpkg)
		if err != nil {
			return nil, fmt.Errorf("G空間情報センターにデータセット %s を作成できませんでした: %w", pkgName, err)
		}
		return &pkg2, nil
	}

	return pkg, nil
}

func (s *handler) findPackage(ctx context.Context, cityCode, cityName string, year int) (_ *ckan.Package, n string, err error) {
	// pattern1 -shi
	name := datasetName(cityCode, cityName, year)
	p, _ := s.ckan.ShowPackage(ctx, name)
	if p.Name != "" {
		return &p, p.Name, nil
	}

	// pattern2 -city
	name2 := datasetName(cityCode, strings.Replace(cityName, "-shi", "-city", 1), year)
	if name != name2 {
		p, _ = s.ckan.ShowPackage(ctx, name2)
		if p.Name != "" {
			return &p, p.Name, nil
		}
	}

	return nil, name, nil
}

type ResourceInfo struct {
	Name        string
	URL         string
	Description string
}

func (resInfo ResourceInfo) Into(pkgID, resID string) ckan.Resource {
	return ckan.Resource{
		ID:        resID,
		PackageID: pkgID,
		Name:      resInfo.Name,
		URL:       resInfo.URL,
	}
}

func (s *handler) createOrUpdateResource(ctx context.Context, pkg *ckan.Package, resInfo ResourceInfo) error {
	// find
	res := findResource(pkg, resInfo.Name)
	if res != nil {
		if _, err := s.ckan.PatchResource(ctx, resInfo.Into(pkg.ID, res.ID)); err != nil {
			return fmt.Errorf("G空間情報センターのリソース %s を更新できませんでした: %w", resInfo.Name, err)
		}

		log.Infofc(ctx, "geospartialjpv3: resource %s updated", resInfo.Name)
		return nil
	}

	_, err := s.ckan.CreateResource(ctx, resInfo.Into(pkg.ID, ""))
	if err != nil {
		return fmt.Errorf("G空間情報センターにリソース %s を作成できませんでした: %w", resInfo.Name, err)
	}

	log.Infofc(ctx, "geospartialjpv3: resource %s created", resInfo.Name)
	return nil
}

func findResource(pkg *ckan.Package, resourceName string) *ckan.Resource {
	res, ok := lo.Find(pkg.Resources, func(r ckan.Resource) bool {
		return r.Name == resourceName
	})
	if !ok {
		return nil
	}
	return &res
}

func datasetName(cityCode, cityName string, year int) string {
	datasetName := ""
	if isTokyo23ku(cityName) {
		if year <= 2020 {
			datasetName = fmt.Sprintf("plateau-%s", gspatialjpTokyo23ku)
		} else {
			datasetName = fmt.Sprintf("plateau-%s-%d", gspatialjpTokyo23ku, year)
		}
	} else {
		datasetName = fmt.Sprintf("plateau-%s-%s-%d", cityCode, cityName, year)
	}
	return datasetName
}

func isTokyo23ku(cityName string) bool {
	return cityName == citygmlTokyo23ku || cityName == citygmlTokyo23ku2 || cityName == gspatialjpTokyo23ku
}

var (
	gspatialjpTokyo23ku = "tokyo23ku"
	citygmlTokyo23ku    = "tokyo23-ku"
	citygmlTokyo23ku2   = "tokyo-23ku"
)
