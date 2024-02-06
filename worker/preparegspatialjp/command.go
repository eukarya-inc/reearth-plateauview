package preparegspatialjp

import (
	"context"
	"fmt"
	"os"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

type Config struct {
	CMSURL     string
	CMSToken   string
	ProjectID  string
	CityItemID string
	WetRun     bool
}

func Command(conf *Config) (err error) {
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
	log.Infofc(ctx, "city item: %s", ppp.Sprint(cityItem))

	cityName := cityItem.CityName
	log.Infofc(ctx, "city name: %s", cityName)

	var citygmlError, plateauError bool
	defer func() {
		var errmsg string
		if err != nil {
			errmsg = err.Error()
		}
		if err := notifyError(ctx, cms, cityItem.GeospatialjpData, citygmlError, plateauError, errmsg); err != nil {
			log.Errorfc(ctx, "failed to notify error: %w", err)
		}
	}()

	log.Infofc(ctx, "getting all feature items...")

	allFeatureItems, err := getAllFeatureItems(ctx, cms, cityItem)
	if err != nil {
		citygmlError = true
		plateauError = true
		return fmt.Errorf("failed to get all feature items: %w", err)
	}

	log.Infofc(ctx, "feature items: %s", ppp.Sprint(allFeatureItems))
	log.Infofc(ctx, "preparing citygml and plateau...")

	type result struct {
		zipName string
		zipPath string
		err     error
	}

	citygmlCh := make(chan result)
	plateauCh := make(chan result)

	if err := notifyRunning(ctx, cms, cityItem.GeospatialjpData, true, true); err != nil {
		return fmt.Errorf("failed to notify running: %w", err)
	}

	go func(c chan result) {
		citygmlZipName, citygmlZipPath, _, _, err := PrepareCityGML(ctx, cms, cityItem, allFeatureItems)
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

	maxlodCh := lo.Async(func() result {
		name, path, err := MergeMaxLOD(ctx, cms, cityItem, allFeatureItems)
		return result{
			zipName: name,
			zipPath: path,
			err:     err,
		}
	})

	citygmlResult := <-citygmlCh
	citygmlZipPath, citygmlZipName, err := citygmlResult.zipPath, citygmlResult.zipName, citygmlResult.err
	if err != nil {
		citygmlError = true
		return fmt.Errorf("failed to prepare citygml: %w", err)
	}

	plateauResult := <-plateauCh
	plateauZipPath, plateauZipName, err := plateauResult.zipPath, plateauResult.zipName, plateauResult.err
	if err != nil {
		plateauError = true
		return fmt.Errorf("failed to prepare plateau: %w", err)
	}

	maxlodResult := <-maxlodCh
	if maxlodResult.err != nil {
		plateauError = true
		return fmt.Errorf("failed to merge maxlod: %w", err)
	}

	var citygmlZipAssetID, plateauZipAssetID, maxlodAssetID string

	relatedZipAssetID, err := GetRelatedZipAssetID(ctx, cms, cityItem)
	if err != nil {
		log.Errorfc(ctx, "failed to get related zip asset id: %w", err)
	}

	if conf.WetRun {
		log.Infofc(ctx, "uploading zips...")

		if citygmlZipAssetID, err = uploadZip(ctx, cms, conf.ProjectID, citygmlZipName, citygmlZipPath); err != nil {
			return fmt.Errorf("failed to upload citygml zip: %w", err)
		}

		if plateauZipAssetID, err = uploadZip(ctx, cms, conf.ProjectID, plateauZipName, plateauZipPath); err != nil {
			return fmt.Errorf("failed to upload plateau zip: %w", err)
		}

		if maxlodAssetID, err = uploadZip(ctx, cms, conf.ProjectID, maxlodResult.zipName, maxlodResult.zipPath); err != nil {
			return fmt.Errorf("failed to upload maxlod: %w", err)
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

	if maxlodAssetID != "" {
		log.Infofc(ctx, "maxlod asset id: %s", maxlodAssetID)
	} else {
		log.Infofc(ctx, "maxlod asset id: (not uploaded)")
	}

	if conf.WetRun {
		log.Infofc(ctx, "attaching assets...")
		if err := attachAssets(ctx, cms, cityItem, citygmlZipAssetID, plateauZipAssetID, relatedZipAssetID, maxlodAssetID); err != nil {
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

func attachAssets(ctx context.Context, c *cms.CMS, cityItem *CityItem, citygmlZipAssetID, plateauZipAssetID, relatedZipAssetID, maxlodAssetID string) error {
	if citygmlZipAssetID == "" && plateauZipAssetID == "" && relatedZipAssetID == "" {
		return nil
	}

	item := GspatialjpItem{
		ID: cityItem.GeospatialjpData,
	}

	if citygmlZipAssetID != "" {
		item.CityGML = citygmlZipAssetID
		item.MergeCityGMLStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if plateauZipAssetID != "" {
		item.Plateau = plateauZipAssetID
		item.MergePlateauStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if relatedZipAssetID != "" {
		item.Related = plateauZipAssetID
		item.MergeRelatedStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if maxlodAssetID != "" {
		item.MaxLOD = maxlodAssetID
	}

	if relatedZipAssetID != "" {
		item.Related = relatedZipAssetID
	}

	var rawItem *cms.Item
	cms.Marshal(item, rawItem)

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	if err := c.CommentToItem(ctx, cityItem.GeospatialjpData, "マージ処理が完了しました。"); err != nil {
		return fmt.Errorf("failed to comment to item: %w", err)
	}

	return nil
}

func notifyError(ctx context.Context, c *cms.CMS, cityItemID string, citygmlError, plateauError bool, comment string) error {
	if comment != "" {
		if err := c.CommentToItem(ctx, cityItemID, fmt.Sprintf("マージ処理に失敗しました。%s", comment)); err != nil {
			return fmt.Errorf("failed to comment to item: %w", err)
		}
	}

	if !citygmlError && !plateauError {
		return nil
	}

	item := GspatialjpItem{
		ID: cityItemID,
	}

	if citygmlError {
		item.MergeCityGMLStatus = &cms.Tag{
			Name: "エラー",
		}
	} else {
		item.MergeCityGMLStatus = &cms.Tag{
			Name: "未実行",
		}
	}

	if plateauError {
		item.MergePlateauStatus = &cms.Tag{
			Name: "エラー",
		}
	} else {
		item.MergePlateauStatus = &cms.Tag{
			Name: "未実行",
		}
	}

	var rawItem *cms.Item
	cms.Marshal(item, rawItem)

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func notifyRunning(ctx context.Context, c *cms.CMS, cityItemID string, citygmlRunning, plateauRunning bool) error {
	if !citygmlRunning && !plateauRunning {
		return nil
	}

	item := GspatialjpItem{
		ID: cityItemID,
	}

	if citygmlRunning {
		item.MergeCityGMLStatus = &cms.Tag{
			Name: "実行中",
		}
	}

	if plateauRunning {
		item.MergePlateauStatus = &cms.Tag{
			Name: "実行中",
		}
	}

	var rawItem cms.Item
	cms.Marshal(item, &rawItem)

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if err := c.CommentToItem(ctx, rawItem.ID, "マージ処理を開始しました。"); err != nil {
		return fmt.Errorf("failed to comment to item: %w", err)
	}

	return nil
}
