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

const tmpDir = "tmp"

type Config struct {
	CMSURL      string
	CMSToken    string
	ProjectID   string
	CityItemID  string
	SkipCityGML bool
	SkipPlateau bool
	SkipMaxLOD  bool
	SkipIndex   bool
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

	// get items fron CNS
	log.Infofc(ctx, "getting item from CMS...")

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

	if cityItem.YearInt() == 0 {
		return fmt.Errorf("invalid year: %s", cityItem.Year)
	}

	if cityItem.SpecVersionMajorInt() == 0 {
		return fmt.Errorf("invalid spec version: %s", cityItem.Spec)
	}

	indexItemRaw, err := cms.GetItem(ctx, cityItem.GeospatialjpIndex, false)
	if err != nil {
		return fmt.Errorf("failed to get index item: %w", err)
	}

	indexItem := GspatialjpIndexItemFrom(indexItemRaw)
	log.Infofc(ctx, "geospatialjp index item: %s", ppp.Sprint(indexItem))

	gdataItemRaw, err := cms.GetItem(ctx, cityItem.GeospatialjpData, false)
	if err != nil {
		return fmt.Errorf("failed to get geospatialjp data item: %w", err)
	}

	gdataItem := GspatialjpDataItemFrom(gdataItemRaw)
	log.Infofc(ctx, "geospatialjp data item: %s", ppp.Sprint(gdataItem))

	if !gdataItem.ShouldMergeCityGML() {
		conf.SkipCityGML = true
	}
	if !gdataItem.ShouldMergePlateau() {
		conf.SkipPlateau = true
	}
	if !gdataItem.ShouldMergeMaxLOD() {
		conf.SkipMaxLOD = true
	}
	if conf.SkipCityGML && conf.SkipPlateau && conf.SkipMaxLOD && conf.SkipRelated {
		return fmt.Errorf("no command to run")
	}

	// do merging
	var comment string
	var citygmlError, plateauError, maxlodError bool
	defer func() {
		if err != nil {
			comment = err.Error()
		}
		if err := notifyError(
			ctx, cms,
			cityItem.GeospatialjpData,
			err != nil,
			citygmlError, plateauError, maxlodError,
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
		maxlodError = true
		return fmt.Errorf("failed to get all feature items: %w", err)
	}

	log.Infofc(ctx, "feature items: %s", ppp.Sprint(allFeatureItems))
	log.Infofc(ctx, "preparing citygml and plateau...")

	if err := notifyRunning(ctx, cms, cityItem.GeospatialjpData, true, true); err != nil {
		return fmt.Errorf("failed to notify running: %w", err)
	}

	citygmlCh := lo.Async(func() lo.Tuple3[string, string, error] {
		if conf.SkipCityGML {
			return lo.Tuple3[string, string, error]{}
		}

		name, path, _, _, err := PrepareCityGML(ctx, cms, cityItem, allFeatureItems)
		if err != nil {
			return lo.Tuple3[string, string, error]{
				C: err,
			}
		}

		return lo.Tuple3[string, string, error]{
			A: name,
			B: path,
		}
	})

	plateauCh := lo.Async(func() lo.Tuple3[string, string, error] {
		if conf.SkipPlateau {
			return lo.Tuple3[string, string, error]{}
		}

		name, path, err := PreparePlateau(ctx, cms, cityItem, allFeatureItems)
		return lo.Tuple3[string, string, error]{
			A: name,
			B: path,
			C: err,
		}
	})

	maxlodCh := lo.Async(func() lo.Tuple3[string, string, error] {
		if conf.SkipMaxLOD {
			return lo.Tuple3[string, string, error]{}
		}

		name, path, err := MergeMaxLOD(ctx, cms, cityItem, allFeatureItems)
		return lo.Tuple3[string, string, error]{
			A: name,
			B: path,
			C: err,
		}
	})

	citygmlResult := <-citygmlCh
	plateauResult := <-plateauCh
	maxlodResult := <-maxlodCh

	// check errors
	if citygmlResult.C != nil || plateauResult.C != nil || maxlodResult.C != nil {
		var errs []error
		if citygmlResult.C != nil {
			citygmlError = true
			errs = append(errs, fmt.Errorf("CityGMLのマージに失敗しました: %w\n", citygmlResult.C))
		}
		if plateauResult.C != nil {
			plateauError = true
			errs = append(errs, fmt.Errorf("3D Tiles,MVTのマージに失敗しました: %w\n", plateauResult.C))
		}
		if maxlodResult.C != nil {
			maxlodError = true
			errs = append(errs, fmt.Errorf("最大LODのマージに失敗しました: %w", maxlodResult.C))
		}
		err = errors.Join(errs...)
		return err
	}

	// generate index
	var index string
	if !conf.SkipIndex {
		index, err = GenerateIndex(ctx, &IndexSeed{
			CityName:       cityItem.CityName,
			Year:           cityItem.YearInt(),
			V:              cityItem.SpecVersionMajorInt(),
			CityGMLZipPath: citygmlResult.B,
			PlateuaZipPath: plateauResult.B,
			RelatedZipPath: "",  //TODO
			Generic:        nil, //TODO
		})
		if err != nil {
			return fmt.Errorf("目録の生成に失敗しました: %w", err)
		}
	}

	var citygmlAssetID, plateauAssetID, maxlodAssetID, relatedAssetID string

	// get related data asset ID
	if !conf.SkipRelated {
		var err2 error
		relatedAssetID, err2 = GetRelatedZipAssetID(ctx, cms, cityItem)
		if err2 != nil {
			log.Errorfc(ctx, "failed to get related zip asset id: %w", err2)
			comment += fmt.Sprintf("\n関連ファイルの取得に失敗しました。: %s", err2)
		}
	}

	// upload zips
	if conf.WetRun {
		log.Infofc(ctx, "uploading zips...")

		if citygmlResult.A != "" && citygmlResult.C == nil {
			if citygmlAssetID, err = upload(ctx, cms, conf.ProjectID, citygmlResult.A, citygmlResult.B); err != nil {
				return fmt.Errorf("failed to upload citygml zip: %w", err)
			}
		}

		if plateauResult.A != "" && plateauResult.C == nil {
			if plateauAssetID, err = upload(ctx, cms, conf.ProjectID, plateauResult.A, plateauResult.B); err != nil {
				return fmt.Errorf("failed to upload plateau zip: %w", err)
			}
		}

		if maxlodResult.A != "" && maxlodResult.C == nil {
			if maxlodAssetID, err = upload(ctx, cms, conf.ProjectID, maxlodResult.A, maxlodResult.B); err != nil {
				return fmt.Errorf("failed to upload maxlod: %w", err)
			}
		}
	}

	// logging
	if citygmlAssetID != "" {
		log.Infofc(ctx, "citygml zip asset id: %s", citygmlAssetID)
	} else {
		log.Infofc(ctx, "citygml zip asset id: (not uploaded)")
	}

	if plateauAssetID != "" {
		log.Infofc(ctx, "plateau zip asset id: %s", plateauAssetID)
	} else {
		log.Infofc(ctx, "plateau zip asset id: (not uploaded)")
	}

	if relatedAssetID != "" {
		log.Infofc(ctx, "related zip asset id: %s", relatedAssetID)
	} else {
		log.Infofc(ctx, "related zip asset id: (not uploaded)")
	}

	if maxlodAssetID != "" {
		log.Infofc(ctx, "maxlod asset id: %s", maxlodAssetID)
	} else {
		log.Infofc(ctx, "maxlod asset id: (not uploaded)")
	}

	// attach assets
	if conf.WetRun && err == nil {
		log.Infofc(ctx, "attaching assets...")
		result := finalResult{
			CityGMLAssetID: citygmlAssetID,
			PlateauAssetID: plateauAssetID,
			RelatedAssetID: relatedAssetID,
			MaxLODAssetID:  maxlodAssetID,
			Index:          index,
		}
		if err := attachAssets(ctx, cms, cityItem, result); err != nil {
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

type finalResult struct {
	CityGMLAssetID string
	PlateauAssetID string
	RelatedAssetID string
	MaxLODAssetID  string
	Index          string
}

func (f finalResult) IsEmpty() bool {
	return f.CityGMLAssetID == "" && f.PlateauAssetID == "" && f.RelatedAssetID == "" && f.MaxLODAssetID == ""
}

func attachAssets(ctx context.Context, c *cms.CMS, cityItem *CityItem, result finalResult) error {
	if cityItem == nil || result.IsEmpty() {
		return nil
	}

	item := GspatialjpDataItem{
		ID: cityItem.GeospatialjpData,
	}

	if result.CityGMLAssetID != "" {
		item.CityGML = result.CityGMLAssetID
		item.MergeCityGMLStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if result.PlateauAssetID != "" {
		item.Plateau = result.PlateauAssetID
		item.MergePlateauStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if result.RelatedAssetID != "" {
		item.Related = result.RelatedAssetID
		item.MergeRelatedStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if result.MaxLODAssetID != "" {
		item.MaxLOD = result.MaxLODAssetID
		item.MergeMaxLODStatus = &cms.Tag{
			Name: "成功",
		}
	}

	if result.Index != "" {
		item.Index = result.Index
	}

	if result.RelatedAssetID != "" {
		item.Related = result.RelatedAssetID
	}

	var rawItem cms.Item
	cms.Marshal(item, &rawItem)

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	if err := c.CommentToItem(ctx, cityItem.GeospatialjpData, "マージ処理が完了しました。"); err != nil {
		return fmt.Errorf("failed to comment to item: %w", err)
	}

	return nil
}

func notifyError(ctx context.Context, c *cms.CMS, cityItemID string, isErr bool, citygmlError, plateauError, maxLODError bool, comment string) error {
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

	item := GspatialjpDataItem{
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

	if maxLODError {
		item.MergeMaxLODStatus = &cms.Tag{
			Name: "エラー",
		}
	} else {
		item.MergeMaxLODStatus = &cms.Tag{
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

	item := GspatialjpDataItem{
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
