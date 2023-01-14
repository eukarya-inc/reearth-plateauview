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
	"github.com/xuri/excelize/v2"
)

var (
	modelKey     = "plateau"
	createEvent  = "item.create"
	updateEvent  = "item.update"
	publishEvent = "item.publish"
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

		if w.Type != createEvent && w.Type != updateEvent && w.Type != publishEvent {
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

		if item.Catalog == "" {
			log.Debugf("geospatialjp webhook: skipped: no catalog")
			return nil
		}

		var err error
		var act string
		if w.Type == publishEvent {
			// publish event: create resources to ckan
			act = "create resources to ckan"
			err = s.CreateCKANResources(ctx, item)
		} else if item.CatalogStatus == "" || item.CatalogStatus == StatusReady {
			// create or update event: check the catalog file
			act = "check catalog"
			err = s.CheckCatalog(ctx, w.Data.Schema.ProjectID, item)

			// comment to item
			comment := fmt.Sprintf("目録ファイルの検査でエラーが発生しました。%s", err)
			if err2 := s.CMS.CommentToItem(ctx, item.ID, comment); err2 != nil {
				log.Errorf("failed to comment to item %s: %s", item.ID, err2)
			}

			// update item
			if _, err2 := s.CMS.UpdateItem(ctx, item.ID, Item{
				CatalogStatus: StatusError,
			}.Fields()); err2 != nil {
				log.Errorf("failed to update item %s: %s", item.ID, err2)
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
	// get catalog url
	catalogAsset, err := s.CMS.Asset(ctx, i.Catalog)
	if err != nil {
		return fmt.Errorf("failed to get catalog asset: %w", err)
	}

	catalogFinalFileName, err := catalogFinalFileName(catalogAsset.URL)
	if err != nil {
		return fmt.Errorf("invalid catalog URL: %w", err)
	}

	catalogAssetRes, err := http.DefaultClient.Do(util.DR(
		http.NewRequestWithContext(ctx, http.MethodGet, catalogAsset.URL, nil)))
	if err != nil {
		return fmt.Errorf("failed to get catalog asset: %w", err)
	}
	if catalogAssetRes.StatusCode != 200 {
		return fmt.Errorf("failed to get catalog asset: status code is %d", catalogAssetRes.StatusCode)
	}

	defer catalogAssetRes.Body.Close()

	// parse
	xf, err := excelize.OpenReader(catalogAssetRes.Body)
	if err != nil {
		return fmt.Errorf("failed to open catalog: %w", err)
	}

	cf := NewCatalogFile(xf)
	c, err := cf.Parse()
	if err != nil {
		if err2 := s.CMS.CommentToItem(ctx, i.ID, fmt.Sprintf("%s", err)); err2 != nil {
			return fmt.Errorf("failed to comment to %s: err = %w, content = %s", i.ID, err2, err)
		}
		return err
	}

	if err := c.Validate(); err != nil {
		if err2 := s.CMS.CommentToItem(ctx, i.ID, fmt.Sprintf("%s", err)); err != nil {
			return fmt.Errorf("failed to comment to %s: err = %w, content = %s", i.ID, err2, err)
		}
		return err
	}

	// delete sheet
	if err := cf.DeleteSheet(); err != nil {
		return fmt.Errorf("failed to delete sheet: %w", err)
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
		return fmt.Errorf("failed to upload catalog: %w", err)
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

func (s *Services) CreateCKANResources(ctx context.Context, i Item) error {
	// upload catalog
	if i.Catalog != "" {
		// get catalog url
		catalogAsset, err := s.CMS.Asset(ctx, i.Catalog)
		if err != nil {
			return fmt.Errorf("failed to get catalog asset: %w", err)
		}

		catalogAssetRes, err := http.DefaultClient.Do(util.DR(
			http.NewRequestWithContext(ctx, http.MethodGet, catalogAsset.URL, nil)))
		if err != nil {
			return fmt.Errorf("failed to get catalog asset: %w", err)
		}
		if catalogAssetRes.StatusCode != 200 {
			return fmt.Errorf("failed to get catalog asset: status code is %d", catalogAssetRes.StatusCode)
		}

		defer catalogAssetRes.Body.Close()

		// parse catalog
		xf, err := excelize.OpenReader(catalogAssetRes.Body)
		if err != nil {
			return fmt.Errorf("failed to open catalog: %w", err)
		}

		c, err := NewCatalogFile(xf).Parse()
		if err != nil {
			if err2 := s.CMS.CommentToItem(ctx, i.ID, fmt.Sprintf("%s", err)); err2 != nil {
				return fmt.Errorf("failed to comment to %s: err = %w, content = %s", i.ID, err2, err)
			}
			return fmt.Errorf("cannot parse catalog: %w", err)
		}

		if err := c.Validate(); err != nil {
			if err2 := s.CMS.CommentToItem(ctx, i.ID, fmt.Sprintf("%s", err)); err != nil {
				return fmt.Errorf("failed to comment to %s: err = %w, content = %s", i.ID, err2, err)
			}
			return fmt.Errorf("invalid catalog: %w", err)
		}
	}

	return nil
}

func (s *Services) UploadCatalogToCkan(ctx context.Context, i Item) error {
	if i.Catalog == "" || i.CatalogFinal == "" {
		return nil
	}

	// get citygml asset
	cityGMLAssetID := i.CityGMLGeoSpatialJP
	if cityGMLAssetID == "" {
		cityGMLAssetID = i.CityGML
	}
	if cityGMLAssetID == "" {
		return errors.New("no citygml")
	}
	citygmlAsset, err := s.CMS.Asset(ctx, cityGMLAssetID)
	if err != nil {
		return fmt.Errorf("failed to get citygml asset: %w", err)
	}
	pkgKey, err := packageKey(citygmlAsset.URL)
	if err != nil {
		return fmt.Errorf("cannot get package key: %w", err)
	}

	//  get all url
	allAsset, err := s.CMS.Asset(ctx, i.All)
	if err != nil {
		return fmt.Errorf("failed to get all asset: %w", err)
	}

	// get catalog url
	catalogAsset, err := s.CMS.Asset(ctx, i.Catalog)
	if err != nil {
		return fmt.Errorf("failed to get catalog asset: %w", err)
	}

	catalogFinalAsset, err := s.CMS.Asset(ctx, i.CatalogFinal)
	if err != nil {
		return fmt.Errorf("failed to get catalog final asset: %w", err)
	}

	// open catalog
	c, err := s.catalog(ctx, catalogAsset.URL)
	if err != nil {
		return err
	}

	// find package
	pkg, err := s.findPackage(ctx, pkgKey, c)
	if err != nil {
		return fmt.Errorf("cannot find package: %w", err)
	}

	resources := resources(pkg, c, citygmlAsset.URL, allAsset.URL, catalogFinalAsset.URL, s.CkanPrivate)

	// register resource
	for _, r := range resources {
		if _, err = s.Ckan.SaveResource(ctx, r); err != nil {
			return fmt.Errorf("failed to register resource: %w", err)
		}
	}

	return nil
}

func (s *Services) catalog(ctx context.Context, catalogURL string) (c Catalog, _ error) {
	catalogAssetRes, err := http.DefaultClient.Do(util.DR(
		http.NewRequestWithContext(ctx, http.MethodGet, catalogURL, nil)))
	if err != nil {
		return c, fmt.Errorf("failed to get catalog asset: %w", err)
	}
	if catalogAssetRes.StatusCode != 200 {
		return c, fmt.Errorf("failed to get catalog asset: status code is %d", catalogAssetRes.StatusCode)
	}

	defer catalogAssetRes.Body.Close()

	// parse catalog
	xf, err := excelize.OpenReader(catalogAssetRes.Body)
	if err != nil {
		return c, fmt.Errorf("failed to open catalog: %w", err)
	}

	cf := NewCatalogFile(xf)
	c, err = cf.Parse()
	if err != nil {
		return c, fmt.Errorf("cannot parse catalog: %w", err)
	}

	return c, nil
}

func (s *Services) findPackage(ctx context.Context, pkgKey string, c Catalog) (*ckan.Package, error) {
	p, err := s.Ckan.ShowPackage(ctx, pkgKey)
	if err != nil {
		p, err = s.Ckan.CreatePackage(ctx, packageFromCatalog(c, pkgKey))
		if err != nil {
			return nil, fmt.Errorf("cannot create package to ckan: %w", err)
		}
	}
	return &p, nil
}
