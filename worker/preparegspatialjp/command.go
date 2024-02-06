package preparegspatialjp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

type Config struct {
	CMSURL      string
	CMSToken    string
	ProjectID   string
	CityItemID  string
	SkipCityGML bool
	SkipPlateau bool
	SkipMaxLOD  bool
	SkipRelated bool
	WetRun      bool
}

func Command(conf *Config) (err error) {
	if conf == nil || conf.SkipCityGML && conf.SkipPlateau && conf.SkipMaxLOD && conf.SkipRelated {
		return fmt.Errorf("no command to run")
	}

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
	log.Infofc(ctx, "city item raw: %s", ppp.Sprint(cityItemRaw))

	cityItem := CityItemFrom(cityItemRaw)
	log.Infofc(ctx, "city item: %s", ppp.Sprint(cityItem))

	if cityItem == nil || cityItem.CityCode == "" || cityItem.CityName == "" || cityItem.CityNameEn == "" || cityItem.GeospatialjpData == "" {
		return fmt.Errorf("invalid city item: %s", conf.CityItemID)
	}

	var comment string
	var citygmlError, plateauError bool
	defer func() {
		if err != nil {
			comment = err.Error()
		}
		if err := notifyError(
			ctx, cms,
			cityItem.GeospatialjpData,
			err != nil,
			citygmlError, plateauError,
			strings.TrimSpace(comment),
		); err != nil {
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
		Name string
		Path string
		Err  error
	}

	if err := notifyRunning(ctx, cms, cityItem.GeospatialjpData, true, true); err != nil {
		return fmt.Errorf("failed to notify running: %w", err)
	}

	citygmlCh := lo.Async(func() result {
		if conf.SkipCityGML {
			return result{}
		}
		name, path, _, _, err := PrepareCityGML(ctx, cms, cityItem, allFeatureItems)
		return result{
			Name: name,
			Path: path,
			Err:  err,
		}
	})

	plateauCh := lo.Async(func() result {
		if conf.SkipPlateau {
			return result{}
		}

		name, path, err := PreparePlateau(ctx, cms, cityItem, allFeatureItems)
		return result{
			Name: name,
			Path: path,
			Err:  err,
		}
	})

	maxlodCh := lo.Async(func() result {
		if conf.SkipMaxLOD {
			return result{}
		}

		name, path, err := MergeMaxLOD(ctx, cms, cityItem, allFeatureItems)
		return result{
			Name: name,
			Path: path,
			Err:  err,
		}
	})

	citygmlResult := <-citygmlCh
	plateauResult := <-plateauCh
	maxlodResult := <-maxlodCh

	if citygmlResult.Err != nil {
		citygmlError = true
	}

	if plateauResult.Err != nil {
		plateauError = true
	}

	if citygmlResult.Err != nil || plateauResult.Err != nil {
		err = errors.Join(citygmlResult.Err, plateauResult.Err)
		return err
	}

	if maxlodResult.Err != nil {
		log.Errorfc(ctx, "failed to merge maxlod: %w", maxlodResult.Err)
		comment += fmt.Sprintf("\n最大LODのマージ処理に失敗しました。: %s", maxlodResult.Err)
	}

	var citygmlZipAssetID, plateauZipAssetID, maxlodAssetID, relatedZipAssetID string

	if !conf.SkipCityGML {
		var err2 error
		relatedZipAssetID, err2 = GetRelatedZipAssetID(ctx, cms, cityItem)
		if err2 != nil {
			log.Errorfc(ctx, "failed to get related zip asset id: %w", err2)
			comment += fmt.Sprintf("\n関連ファイルの取得に失敗しました。: %s", err2)
		}
	}

	if conf.WetRun {
		log.Infofc(ctx, "uploading zips...")

		if citygmlResult.Name != "" && citygmlResult.Err == nil {
			if citygmlZipAssetID, err = upload(ctx, cms, conf.ProjectID, citygmlResult.Name, citygmlResult.Path); err != nil {
				return fmt.Errorf("failed to upload citygml zip: %w", err)
			}
		}

		if plateauResult.Name != "" && plateauResult.Err == nil {
			if plateauZipAssetID, err = upload(ctx, cms, conf.ProjectID, plateauResult.Name, plateauResult.Path); err != nil {
				return fmt.Errorf("failed to upload plateau zip: %w", err)
			}
		}

		if maxlodResult.Name != "" && maxlodResult.Err == nil {
			if maxlodAssetID, err = upload(ctx, cms, conf.ProjectID, maxlodResult.Name, maxlodResult.Path); err != nil {
				return fmt.Errorf("failed to upload maxlod: %w", err)
			}
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

	if conf.WetRun && err == nil {
		log.Infofc(ctx, "attaching assets...")
		if err := attachAssets(ctx, cms, cityItem, citygmlZipAssetID, plateauZipAssetID, relatedZipAssetID, maxlodAssetID); err != nil {
			return fmt.Errorf("failed to attach assets: %w", err)
		}
	}

	return
}

func upload(ctx context.Context, cms *cms.CMS, project, name, path string) (string, error) {
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

func notifyError(ctx context.Context, c *cms.CMS, cityItemID string, isErr bool, citygmlError, plateauError bool, comment string) error {
	if comment != "" {
		msgPrefix := ""
		if isErr {
			msgPrefix = "マージ処理に失敗しました。"
		} else {
			msgPrefix = "マージ処理が完了しました。"
		}
		if err := c.CommentToItem(ctx, cityItemID, msgPrefix+comment); err != nil {
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

	var rawItem cms.Item
	cms.Marshal(item, &rawItem)
	if rawItem.ID == "" {
		return fmt.Errorf("failed to marshal item")
	}

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
