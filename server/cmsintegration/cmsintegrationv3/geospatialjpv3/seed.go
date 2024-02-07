package geospatialjpv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
)

type PackageName struct {
	CityCode, CityNameEn string
	Year                 int
}

func (p PackageName) String() string {
	return datasetName(p.CityCode, p.CityNameEn, p.Year)
}

func PackageNameFrom(city *CityItem) PackageName {
	return PackageName{
		CityCode:   city.CityCode,
		CityNameEn: city.CityNameEn,
		Year:       city.YearInt(),
	}
}

type PackageSeed struct {
	Name            PackageName
	NameJa          string
	Description     string
	OwnerOrg        string
	Area            string
	ThumbnailURL    string
	Author          string
	AuthorEmail     string
	Maintainer      string
	MaintainerEmail string
	Quality         string
	Version         string
}

func (p PackageSeed) Title() string {
	return fmt.Sprintf("3D都市モデル（Project PLATEAU）%s（%d年度）", p.NameJa, p.Name.Year)
}

func (p PackageSeed) ToPackage() ckan.Package {
	return ckan.Package{
		Notes:        p.Description,
		ThumbnailURL: p.ThumbnailURL,
	}
}

func (p PackageSeed) ToNewPackage() ckan.Package {
	tags := append([]ckan.Tag{}, defaultTags...)
	tags = append(tags, ckan.Tag{
		Name: p.NameJa,
	})

	return ckan.Package{
		Name:             p.Name.String(),
		Title:            p.Title(),
		OwnerOrg:         p.OwnerOrg,
		Notes:            p.Description,
		Private:          true,
		Area:             p.Area,
		ThumbnailURL:     p.ThumbnailURL,
		URL:              urlDefault,
		LicenseID:        licenseDefaultID,
		LicenseTitle:     licenseDefaultTitle,
		LicenseURL:       licenseDefaultURL,
		Restriction:      restriction,
		LicenseAgreement: licenseAgreement,
		Fee:              fee,
		Emergency:        emergency,
		Tags:             tags,
		Version:          p.Version,
	}
}

func PackageSeedFrom(cityItem *CityItem, seed Seed) PackageSeed {
	return PackageSeed{
		Name:            PackageNameFrom(cityItem),
		NameJa:          cityItem.CityName,
		OwnerOrg:        seed.Org,
		Description:     seed.Desc,
		Area:            seed.Area,
		ThumbnailURL:    seed.ThumbnailURL,
		Author:          seed.Author,
		AuthorEmail:     seed.AuthorEmail,
		Maintainer:      seed.Maintainer,
		MaintainerEmail: seed.MaintainerEmail,
		Quality:         seed.Quality,
		Version:         cityItem.SpecVersion(),
	}
}
