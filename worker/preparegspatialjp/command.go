package preparegspatialjp

import (
	"context"
	"fmt"
	"os"

	"github.com/k0kubun/pp/v3"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type Config struct {
	CMSURL     string
	CMSToken   string
	ProjectID  string
	CityItemID string
	WetRun     bool
}

func Command(conf *Config) error {
	ctx := context.Background()

	cms, err := cms.New(conf.CMSURL, conf.CMSToken)
	if err != nil {
		return fmt.Errorf("failed to initialize CMS client: %w", err)
	}

	log.Infofc(ctx, "getting city item...")

	cityItemRaw, err := cms.GetItem(ctx, conf.CityItemID, false)
	if err != nil {
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	{
		pp := pp.New()
		pp.SetColoringEnabled(false)
		s := pp.Sprint(cityItem)
		log.Infofc(ctx, "city item: %s", s)
	}

	cityName := cityItem.CityName
	log.Infofc(ctx, "city name: %s", cityName)

	log.Infofc(ctx, "getting all feature items...")

	allFeatureItems, err := getAllFeatureItems(ctx, cms, cityItem)
	if err != nil {
		return fmt.Errorf("failed to get all feature items: %w", err)
	}

	{
		pp := pp.New()
		pp.SetColoringEnabled(false)
		s := pp.Sprint(allFeatureItems)
		log.Infofc(ctx, "feature items: %s", s)
	}

	log.Infofc(ctx, "preparing citygml and plateau...")

	type result struct {
		zipName string
		zipPath string
		err     error
	}

	citygmlCh := make(chan result)
	plateauCh := make(chan result)

	go func(c chan result) {
		citygmlZipName, citygmlZipPath, err := PrepareCityGML(ctx, cms, cityItem, allFeatureItems)
		c <- result{
			zipName: citygmlZipName,
			zipPath: citygmlZipPath,
			err:     err,
		}
	}(citygmlCh)

	go func(c chan result) {
		plateauZipName, plateauZipPath, err := PreparePlateau(ctx, cms, cityItem, allFeatureItems)
		c <- result{
			zipName: plateauZipName,
			zipPath: plateauZipPath,
			err:     err,
		}
	}(plateauCh)

	citygmlResult := <-citygmlCh
	citygmlZipPath, citygmlZipName, err := citygmlResult.zipPath, citygmlResult.zipName, citygmlResult.err
	if err != nil {
		return fmt.Errorf("failed to prepare citygml: %w", err)
	}

	plateauResult := <-plateauCh
	plateauZipPath, plateauZipName, err := plateauResult.zipPath, plateauResult.zipName, plateauResult.err
	if err != nil {
		return fmt.Errorf("failed to prepare plateau: %w", err)
	}

	var citygmlZipAssetID, plateauZipAssetID string

	relatedZipAssetID, err := GetRelatedZipAssetID(ctx, cms, cityItem)
	if err != nil {
		return fmt.Errorf("failed to get related zip asset id: %w", err)
	}

	if conf.WetRun {
		log.Infofc(ctx, "uploading zips...")

		if citygmlZipAssetID, err = uploadZip(ctx, cms, conf.ProjectID, citygmlZipName, citygmlZipPath); err != nil {
			return fmt.Errorf("failed to upload citygml zip: %w", err)
		}

		if plateauZipAssetID, err = uploadZip(ctx, cms, conf.ProjectID, plateauZipName, plateauZipPath); err != nil {
			return fmt.Errorf("failed to upload plateau zip: %w", err)
		}
	}

	if citygmlZipAssetID != "" {
		log.Infofc(ctx, "citygml zip asset id: %s", citygmlZipAssetID)
	} else {
		log.Infofc(ctx, "citygml zip asset id: (not uploaded)")
	}

	if plateauZipAssetID != "" {
		log.Infofc(ctx, "plateau zip asset id: %s", plateauZipAssetID)
	} else {
		log.Infofc(ctx, "plateau zip asset id: (not uploaded)")
	}

	if relatedZipAssetID != "" {
		log.Infofc(ctx, "related zip asset id: %s", relatedZipAssetID)
	} else {
		log.Infofc(ctx, "related zip asset id: (not uploaded)")
	}

	if conf.WetRun {
		log.Infofc(ctx, "attaching assets...")
		if err := attachAssets(ctx, cms, cityItem, citygmlZipAssetID, plateauZipAssetID, relatedZipAssetID); err != nil {
			return fmt.Errorf("failed to attach assets: %w", err)
		}
	}

	return nil
}

func uploadZip(ctx context.Context, cms *cms.CMS, project, name, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	assetID, err := cms.UploadAssetDirectly(ctx, project, name, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload asset: %w", err)
	}

	return assetID, nil
}

func attachAssets(ctx context.Context, c *cms.CMS, cityItem *CityItem, citygmlZipAssetID, plateauZipAssetID, relatedZipAssetID string) error {
	if citygmlZipAssetID == "" && plateauZipAssetID == "" && relatedZipAssetID == "" {
		return nil
	}

	var fields []*cms.Field

	if citygmlZipAssetID != "" {
		fields = append(fields, &cms.Field{
			Key:   "citygml",
			Value: citygmlZipAssetID,
		})
	}

	if plateauZipAssetID != "" {
		fields = append(fields, &cms.Field{
			Key:   "plateau",
			Value: plateauZipAssetID,
		})
	}

	if relatedZipAssetID != "" {
		fields = append(fields, &cms.Field{
			Key:   "related",
			Value: relatedZipAssetID,
		})
	}

	if _, err := c.UpdateItem(ctx, cityItem.GeospatialjpData, fields, nil); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}
