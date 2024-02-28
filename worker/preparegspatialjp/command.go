package preparegspatialjp

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

const tmpDirBase = "plateau-api-worker-tmp"

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

	codelists, err := cms.Asset(ctx, cityItem.CodeLists)
	if err != nil {
		return fmt.Errorf("failed to get codelists: %w", err)
	}

	uc := UpdateCount(codelists.URL)
	if uc == 0 {
		return fmt.Errorf("invalid update count: %s", codelists.URL)
	}

	indexItemRaw, err := cms.GetItem(ctx, cityItem.GeospatialjpIndex, false)
	if err != nil {
		return fmt.Errorf("failed to get index item: %w", err)
	}

	indexItem := GspatialjpIndexItemFrom(indexItemRaw)
	log.Infofc(ctx, "geospatialjp index item: %s", ppp.Sprint(indexItem))

	relatedAssetID, relatedAssetURL, err := GetRelatedZipAssetIDAndURL(ctx, cms, cityItem)
	if err != nil {
		return fmt.Errorf("failed to get related zip asset id and url: %w", err)
	}

	gdataItemRaw, err := cms.GetItem(ctx, cityItem.GeospatialjpData, false)
	if err != nil {
		return fmt.Errorf("failed to get geospatialjp data item: %w", err)
	}

	gdataItem := GspatialjpDataItemFrom(gdataItemRaw)
	log.Infofc(ctx, "geospatialjp data item: %s", ppp.Sprint(gdataItem))

	if gdataItem != nil {
		if !gdataItem.ShouldMergeCityGML() {
			conf.SkipCityGML = true
		}
		if !gdataItem.ShouldMergePlateau() {
			conf.SkipPlateau = true
		}
		if !gdataItem.ShouldMergeMaxLOD() {
			conf.SkipMaxLOD = true
		}
	}

	if conf.SkipCityGML && conf.SkipPlateau && conf.SkipMaxLOD && conf.SkipRelated {
		return fmt.Errorf("no command to run")
	}

	tmpDirName := fmt.Sprintf("%s-%d", time.Now().Format("20060102-150405"), rand.Intn(1000))
	tmpDir := filepath.Join(tmpDirBase, tmpDirName)
	log.Infofc(ctx, "tmp dir: %s", tmpDir)

	// do merging
	var comment string
	var citygmlError, plateauError, maxlodError bool
	defer func() {
		if !conf.WetRun {
			return
		}
		if err != nil {
			comment = err.Error()
		}
		if err := notifyError(
			ctx, cms,
			cityItem.ID,
			cityItem.GeospatialjpData,
			err != nil,
			citygmlError, plateauError, maxlodError,
			strings.TrimSpace(comment),
		); err != nil {
			log.Errorfc(ctx, "failed to notify error: %v", err)
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

	dic := mergeDics(lo.MapToSlice(allFeatureItems, func(k string, v FeatureItem) string {
		return v.Dic
	})...)

	log.Infofc(ctx, "feature items: %s", ppp.Sprint(allFeatureItems))
	log.Infofc(ctx, "dic: %s", ppp.Sprint(dic))
	log.Infofc(ctx, "preparing citygml and plateau...")

	if conf.WetRun {
		if err := notifyRunning(ctx, cms, cityItem.ID, cityItem.GeospatialjpData, !conf.SkipCityGML, !conf.SkipPlateau, !conf.SkipRelated, !conf.SkipMaxLOD); err != nil {
			return fmt.Errorf("failed to notify running: %w", err)
		}
	}

	citygmlCh := lo.Async(func() lo.Tuple3[string, string, error] {
		if conf.SkipCityGML {
			return lo.Tuple3[string, string, error]{}
		}

		name, path, err := PrepareCityGML(ctx, cms, tmpDir, cityItem, allFeatureItems, uc)
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

		name, path, err := PreparePlateau(ctx, cms, tmpDir, cityItem, allFeatureItems, uc)
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

		name, path, err := MergeMaxLOD(ctx, cms, tmpDir, cityItem, allFeatureItems)
		return lo.Tuple3[string, string, error]{
			A: name,
			B: path,
			C: err,
		}
	})

	relatedCh := lo.Async(func() lo.Tuple2[string, error] {
		if conf.SkipRelated || relatedAssetURL == "" {
			return lo.Tuple2[string, error]{}
		}

		p, err := downloadFileTo(ctx, relatedAssetURL, tmpDir)
		return lo.Tuple2[string, error]{
			A: p,
			B: err,
		}
	})

	citygmlResult := <-citygmlCh
	plateauResult := <-plateauCh
	maxlodResult := <-maxlodCh
	relatedResult := <-relatedCh

	// check errors
	if citygmlResult.C != nil || plateauResult.C != nil || maxlodResult.C != nil || relatedResult.B != nil {
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
		if relatedResult.B != nil {
			errs = append(errs, fmt.Errorf("関連ファイルのダウンロードに失敗しました: %w", relatedResult.B))
		}

		err = errors.Join(errs...)
		return err
	}

	// generate index
	var index string
	if !conf.SkipIndex {
		index, err = GenerateIndex(ctx, &IndexSeed{
			CityName:       cityItem.CityName,
			CityCode:       cityItem.CityCode,
			Year:           cityItem.YearInt(),
			V:              cityItem.SpecVersionMajorInt(),
			CityGMLZipPath: citygmlResult.B,
			PlateuaZipPath: plateauResult.B,
			RelatedZipPath: relatedResult.A,
			Generic:        indexItem.Generic,
			Dic:            dic,
		})
		if err != nil {
			return fmt.Errorf("目録の生成に失敗しました: %w", err)
		}
	}

	var citygmlAssetID, plateauAssetID, maxlodAssetID string

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

func notifyError(ctx context.Context, c *cms.CMS, cityItemID, gdataItemID string, isErr bool, citygmlError, plateauError, maxLODError bool, comment string) error {
	if comment != "" {
		msgPrefix := ""
		if isErr {
			msgPrefix = "公開準備処理に失敗しました。"
		} else {
			msgPrefix = "公開準備処理が完了しました。"
		}

		if err := c.CommentToItem(ctx, cityItemID, msgPrefix+comment); err != nil {
			return fmt.Errorf("failed to comment to citygml item: %w", err)
		}

		if err := c.CommentToItem(ctx, gdataItemID, msgPrefix+comment); err != nil {
			return fmt.Errorf("failed to comment to data item: %w", err)
		}
	}

	if !citygmlError && !plateauError && !maxLODError {
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

	if rawItem.Fields == nil {
		rawItem.Fields = []*cms.Field{}
	}

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func notifyRunning(ctx context.Context, c *cms.CMS, citygmlID, dataID string, citygmlRunning, plateauRunning, relatedRunning, maxlodRunning bool) error {
	if !citygmlRunning && !plateauRunning && !relatedRunning && !maxlodRunning {
		return nil
	}

	item := GspatialjpDataItem{
		ID: dataID,
	}

	if citygmlRunning {
		item.MergeCityGMLStatus = &cms.Tag{
			Name: running,
		}
	}

	if plateauRunning {
		item.MergePlateauStatus = &cms.Tag{
			Name: running,
		}
	}

	if relatedRunning {
		item.MergeRelatedStatus = &cms.Tag{
			Name: running,
		}
	}

	if maxlodRunning {
		item.MergeMaxLODStatus = &cms.Tag{
			Name: running,
		}
	}

	var rawItem cms.Item
	cms.Marshal(item, &rawItem)

	if rawItem.Fields == nil {
		rawItem.Fields = []*cms.Field{}
	}

	if _, err := c.UpdateItem(ctx, rawItem.ID, rawItem.Fields, rawItem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	comment := "G空間情報センターの公開準備処理を開始しました。"

	if err := c.CommentToItem(ctx, citygmlID, comment); err != nil {
		return fmt.Errorf("failed to comment to city item: %w", err)
	}

	if err := c.CommentToItem(ctx, dataID, comment); err != nil {
		return fmt.Errorf("failed to comment to data item: %w", err)
	}

	return nil
}
