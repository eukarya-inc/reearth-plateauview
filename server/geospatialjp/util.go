package geospatialjp

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/geospatialjp/ckan"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/vincent-petithory/dataurl"
)

const (
	ResourceNameCityGML = "CityGML（v2）"
	ResourceNameAll     = "3D Tiles, MVT（v2）"
	ResourceNameCatalog = "データ目録（v2）"
	licenseID           = "license_id"
	licenseTitle        = "PLATEAU Site Policy 「３．著作権について」に拠る"
	licenseURL          = "https://www.mlit.go.jp/plateau/site-policy/"
)

func catalogFinalFileName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	b := path.Base(u.Path)

	be, af, found := strings.Cut(b, ".")
	if !found {
		return "", fmt.Errorf("invalid file name: %s", b)
	}

	return fmt.Sprintf("%s_final.%s", be, af), nil
}

func resources(pkg *ckan.Package, c Catalog, citygmlURL, allURL, catalogURL string, private bool) (res []ckan.Resource) {
	citygml, found := lo.Find(pkg.Resources, func(r ckan.Resource) bool {
		return r.Name == ResourceNameCityGML
	})
	if !found {
		citygml = ckan.Resource{
			PackageID: pkg.ID,
			Name:      ResourceNameCityGML,
			Format:    "ZIP",
		}
	}
	if citygml.URL != citygmlURL {
		citygml.URL = citygmlURL
		res = append(res, citygml)
	}

	all, found := lo.Find(pkg.Resources, func(r ckan.Resource) bool {
		return r.Name == ResourceNameAll
	})
	if !found {
		all = ckan.Resource{
			PackageID: pkg.ID,
			Name:      ResourceNameAll,
			Format:    "ZIP",
		}
	}
	if all.URL != allURL {
		all.URL = allURL
		res = append(res, all)
	}

	catalog, found := lo.Find(pkg.Resources, func(r ckan.Resource) bool {
		return r.Name == ResourceNameCatalog
	})
	if !found {
		catalog = ckan.Resource{
			PackageID: pkg.ID,
			Name:      ResourceNameCatalog,
			Format:    "XLSX",
		}
	}
	if catalog.URL != catalogURL {
		catalog.URL = catalogURL
		res = append(res, catalog)
	}

	return res
}

var reFileName = regexp.MustCompile(`^([0-9]+?)_(.+?)_`)

func extractCityName(fn string) (string, string, error) {
	u, err := url.Parse(fn)
	if err != nil {
		return "", "", err
	}

	base := path.Base(u.Path)
	s := reFileName.FindStringSubmatch(base)
	if s == nil {
		return "", "", errors.New("invalid file name")
	}

	return s[1], s[2], nil
}

func packageFromCatalog(c Catalog, org, pkgName string, private bool) ckan.Package {
	// extras := map[string]string{}
	// if c.DataQuality != "" {
	// 	extras["データ品質"] = c.DataQuality
	// }
	// if c.Restriction != "" {
	// 	extras["制約"] = c.Restriction
	// }
	// if c.DisasterClassification != "" {
	// 	extras["災害時区分"] = c.DisasterClassification
	// }
	// if c.FreeOrProvidedClassification != "" {
	// 	extras["有償無償区分"] = c.FreeOrProvidedClassification
	// }
	// if c.GeoArea != "" {
	// 	extras["地理的範囲"] = c.GeoArea
	// }
	// if c.License != "" {
	// 	extras["ライセンス"] = c.License
	// }
	// if c.LicenseAgreement != "" {
	// 	extras["使用許諾"] = c.LicenseAgreement
	// }
	// if c.Organization != "" {
	// 	extras["組織"] = c.Organization
	// }
	// if c.Price != "" {
	// 	extras["価格情報"] = c.Price
	// }
	// if c.Public != "" {
	// 	extras["公開・非公開"] = c.Public
	// }
	// if c.RegisteredDate != "" {
	// 	extras["データ登録日"] = c.RegisteredDate
	// }
	// if c.Source != "" {
	// 	extras["ソース"] = c.Source
	// }

	var thumbnailURL string
	if c.Thumbnail != nil {
		thumbnailURL = dataurl.New(c.Thumbnail, http.DetectContentType(c.Thumbnail)).String()
	}

	return ckan.Package{
		Name:            pkgName,
		Title:           c.Title,
		Private:         private || c.Public != "パブリック",
		Author:          c.Author,
		AuthorEmail:     c.AuthorEmail,
		Maintainer:      c.Maintainer,
		MaintainerEmail: c.MaintainerEmail,
		Notes:           c.Notes,
		Version:         c.Version,
		// Tags:            c.Tags,
		OwnerOrg:         org,
		Restriction:      c.Restriction,
		Charge:           c.Charge,
		RegisteredDate:   c.RegisteredDate,
		LicenseAgreement: c.LicenseAgreement,
		LicenseTitle:     licenseTitle,
		LicenseURL:       licenseURL,
		Fee:              c.Fee,
		Area:             c.Area,
		Quality:          c.Quality,
		Emergency:        c.Emergency,
		URL:              c.Source,
		LicenseID:        licenseID,
		ThumbnailURL:     thumbnailURL,
		// unused:
		// URL: c.URL
		// ライセンス: c.License
		// 組織: c.Organization
		// spatial*: c.Spatial
	}
}
