package geospatialjp

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cms/cmswebhook"
	"github.com/eukarya-inc/reearth-plateauview/server/geospatialjp/ckan"
	"github.com/pkg/errors"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

var (
	modelKey    = "plateau"
	initialYear = 2020
)

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	s, err := NewServices(conf)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, w *cmswebhook.Payload) error {
		if !w.Operator.IsUser() {
			log.Debugf("geospatialjp webhook: invalid event operator: %+v", w.Operator)
			return nil
		}

		if w.Type != cmswebhook.EventItemCreate && w.Type != cmswebhook.EventItemUpdate && w.Type != cmswebhook.EventItemPublish {
			log.Debugf("geospatialjp webhook: invalid event type: %s", w.Type)
			return nil
		}

		if w.Data.Item == nil || w.Data.Model == nil {
			log.Debugf("geospatialjp webhook: invalid event data: %+v", w.Data)
			return nil
		}

		if w.Data.Model.Key != modelKey {
			log.Debugf("geospatialjp webhook: invalid model id: %s, key: %s", w.Data.Item.ModelID, w.Data.Model.Key)
			return nil
		}

		ctx := req.Context()
		item := ItemFrom(*w.Data.Item)

		var err error
		var act string
		if w.Type == cmswebhook.EventItemPublish {
			// publish event: create resources to ckan
			act = "create resources to ckan"
			err = s.RegisterCkanResources(ctx, item)

			if err != nil {
				comment := fmt.Sprintf("G空間情報センターへの登録処理でエラーが発生しました。%s", err)
				s.commentToItem(ctx, item.ID, comment)
			} else {
				s.commentToItem(ctx, item.ID, "G空間情報センターへの登録が完了しました")
			}
		} else {
			// create or update event: check the catalog file
			act = "check catalog"
			err = s.CheckCatalog(ctx, w.Data.Schema.ProjectID, item)

			if err != nil {
				comment := fmt.Sprintf("目録ファイルの検査でエラーが発生しました。%s", err)
				s.commentToItem(ctx, item.ID, comment)

				// update item
				if _, err2 := s.CMS.UpdateItem(ctx, item.ID, Item{
					CatalogStatus: StatusError,
				}.Fields()); err2 != nil {
					log.Errorf("failed to update item %s: %s", item.ID, err2)
				}
			} else {
				s.commentToItem(ctx, item.ID, "目録ファイルの検査とG空間情報センター用目録ファイルの登録が完了しました。")
			}
		}

		if err != nil {
			log.Errorf("geospatialjp webhook: failed to %s: %s", act, err)
		}

		log.Infof("geospatialjp webhook: done")
		return nil
	}, nil
}

func (s *Services) CheckCatalog(ctx context.Context, projectID string, i Item) error {
	if i.CatalogStatus != "" && i.CatalogStatus != StatusReady {
		return nil
	}

	// update item
	if _, err := s.CMS.UpdateItem(ctx, i.ID, Item{
		CatalogStatus: StatusProcessing,
	}.Fields()); err != nil {
		return fmt.Errorf("failed to update item %s: %w", i.ID, err)
	}

	// get catalog url
	catalogAsset, err := s.CMS.Asset(ctx, i.Catalog)
	if err != nil {
		return fmt.Errorf("目録アセットの読み込みに失敗しました。該当アセットが削除されていませんか？: %w", err)
	}

	catalogFinalFileName, err := catalogFinalFileName(catalogAsset.URL)
	if err != nil {
		return fmt.Errorf("invalid catalog URL: %w", err)
	}

	// parse catalog
	c, cf, err := s.parseCatalog(ctx, catalogAsset.URL)
	if err != nil {
		return err
	}

	// validate catalog
	if err := c.Validate(); err != nil {
		return err
	}

	// delete sheet
	if err := cf.DeleteSheet(); err != nil {
		return fmt.Errorf("目録内のG空間情報センター用メタデータシートの削除に失敗しました。: %w", err)
	}

	// upload catalog
	pr, pw := io.Pipe()

	go func() {
		var err error
		defer func() {
			_ = pw.CloseWithError(err)
		}()
		_, err = cf.File().WriteTo(pw)
	}()

	catalogFinalAsset, err := s.CMS.UploadAssetDirectly(ctx, projectID, catalogFinalFileName, pr)
	if err != nil {
		return fmt.Errorf("G空間情報センター用目録のアップロードに失敗しました。: %w", err)
	}

	// update item
	if _, err := s.CMS.UpdateItem(ctx, i.ID, Item{
		CatalogStatus: StatusOK,
		CatalogFinal:  catalogFinalAsset,
	}.Fields()); err != nil {
		return fmt.Errorf("failed to update item %s: %w", i.ID, err)
	}

	return nil
}

func (s *Services) RegisterCkanResources(ctx context.Context, i Item) error {
	if i.Catalog == "" || i.CatalogFinal == "" {
		return errors.New("目録ファイルが登録されていません。目録の検査が正常に完了しているか確認してください。")
	}

	if i.All == "" {
		return errors.New("全データファイルが登録されていません。CityGMLの3D Tiles等への変換が正常に完了しているか確認してください。")
	}

	// get citygml asset
	cityGMLAssetID := i.CityGMLGeoSpatialJP
	if cityGMLAssetID == "" {
		cityGMLAssetID = i.CityGML
	}
	if cityGMLAssetID == "" {
		return errors.New("CityGMLデータが登録されていません。")
	}

	citygmlAsset, err := s.CMS.Asset(ctx, cityGMLAssetID)
	if err != nil {
		return fmt.Errorf("CityGMLデータの読み込みに失敗しました。該当アセットが削除されていませんか？: %w", err)
	}

	cityCode, cityName, err := extractCityName(citygmlAsset.URL)
	if err != nil {
		return fmt.Errorf("CityGMLのzipファイル名から市区町村コードまたは市区町村英名を読み取ることができませんでした。ファイル名の形式が正しいか確認してください。: %w", err)
	}

	// get all url
	allAsset, err := s.CMS.Asset(ctx, i.All)
	if err != nil {
		return fmt.Errorf("全データのアセットの読み込みに失敗しました。該当アセットが削除されていませんか？: %w", err)
	}

	// get catalog url
	catalogAsset, err := s.CMS.Asset(ctx, i.Catalog)
	if err != nil {
		return fmt.Errorf("目録アセットの読み込みに失敗しました。該当アセットが削除されていませんか？: %w", err)
	}

	catalogFinalAsset, err := s.CMS.Asset(ctx, i.CatalogFinal)
	if err != nil {
		return fmt.Errorf("G空間情報センター用の目録アセットの読み込みに失敗しました。該当アセットが削除されていませんか？: %w", err)
	}

	// open catalog
	c, _, err := s.parseCatalog(ctx, catalogAsset.URL)
	if err != nil {
		return err
	}

	// find or create package
	pkg, err := s.findOrCreatePackage(ctx, c, cityCode, cityName)
	if err != nil {
		return err
	}

	// register resource
	resources := resources(pkg, c, citygmlAsset.URL, allAsset.URL, catalogFinalAsset.URL, s.CkanPrivate)
	for _, r := range resources {
		if _, err = s.Ckan.SaveResource(ctx, r); err != nil {
			return fmt.Errorf("G空間情報センターへのリソースの登録に失敗しました。: %w", err)
		}
	}

	return nil
}

func (s *Services) parseCatalog(ctx context.Context, catalogURL string) (c Catalog, cf *CatalogFile, _ error) {
	catalogAssetRes, err := http.DefaultClient.Do(util.DR(
		http.NewRequestWithContext(ctx, http.MethodGet, catalogURL, nil)))
	if err != nil {
		return c, cf, fmt.Errorf("アセットの取得に失敗しました: %w", err)
	}
	if catalogAssetRes.StatusCode != 200 {
		return c, cf, fmt.Errorf("アセットの取得に失敗しました: ステータスコード %d", catalogAssetRes.StatusCode)
	}

	defer catalogAssetRes.Body.Close()

	// parse catalog
	xf, err := excelize.OpenReader(catalogAssetRes.Body)
	if err != nil {
		return c, cf, fmt.Errorf("目録を開くことできませんでした: %w", err)
	}

	cf = NewCatalogFile(xf)
	c, err = cf.Parse()
	if err != nil {
		return c, cf, fmt.Errorf("目録の読み込みに失敗しました: %w", err)
	}

	return c, cf, nil
}

func (s *Services) findOrCreatePackage(ctx context.Context, c Catalog, cityCode, cityName string) (*ckan.Package, error) {
	pkg, pkgName := s.findPackage(ctx, cityCode, cityName)

	if pkg == nil {
		log.Infof("geospartialjp: package plateau-%s-%s-202x not found", cityCode, cityName)

		pkg = lo.ToPtr(packageFromCatalog(c, s.CkanOrg, pkgName, s.CkanPrivate))
		pkg2, err := s.Ckan.CreatePackage(ctx, *pkg)
		if err != nil {
			return nil, fmt.Errorf("G空間情報センターにデータセット %s を作成することができませんでした: %w", pkgName, err)
		}
		pkg = &pkg2
	}
	return pkg, nil
}

func (s *Services) findPackage(ctx context.Context, cityCode, cityName string) (_ *ckan.Package, n string) {
	if cityName == "tokyo23ku" {
		pkg, err := s.Ckan.ShowPackage(ctx, "plateau-tokyo23ku")
		if err != nil {
			return nil, ""
		}
		return &pkg, "plateau-tokyo23ku"
	}

	currentYear := util.Now().Year()
	for y := initialYear; y <= currentYear; y++ {
		n = fmt.Sprintf("plateau-%s-%s-%d", cityCode, cityName, y)
		p, err := s.Ckan.ShowPackage(ctx, n)
		if err == nil {
			return &p, n
		}
	}

	return nil, n
}

func (s *Services) commentToItem(ctx context.Context, itemID, comment string) {
	if err2 := s.CMS.CommentToItem(ctx, itemID, comment); err2 != nil {
		log.Errorf("failed to comment to item %s: %s", itemID, err2)
	}
}
