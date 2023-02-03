package searchindex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/cms/cmswebhook"
	"github.com/reearth/reearthx/log"
)

var (
	modelKey   = "plateau"
	errSkipped = errors.New("not decompressed")
)

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	c, err := cms.New(conf.CMSBase, conf.CMSToken)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, w *cmswebhook.Payload) error {
		if w.Type != cmswebhook.EventItemCreate && w.Type != cmswebhook.EventItemUpdate && w.Type != cmswebhook.EventAssetDecompress {
			log.Debugf("searchindex webhook: invalid event type: %s", w.Type)
			return nil
		}

		pid := w.ProjectID()
		if pid == "" {
			log.Debugf("searchindex webhook: invalid payload: no project id")
			return nil
		}

		ctx := req.Context()

		stprj := conf.CMSStorageProject
		if stprj == "" {
			stprj = pid
		}
		st := NewStorage(c, stprj, conf.CMSStorageModel)

		item, siid, err := getItem(ctx, c, st, w)
		if err != nil || item.ID == "" {
			log.Errorf("searchindex webhook: failed to get item: %v", err)
			return nil
		}

		log.Infof("searchindex webhook: item: %+v", item)

		if item.SearchIndexStatus != "" && item.SearchIndexStatus != StatusReady {
			log.Infof("searchindex webhook: skipped: %s", item.SearchIndexStatus)
			return nil
		}

		if len(item.Bldg) == 0 {
			log.Infof("searchindex webhook: skipped: no bldg assets")
			return nil
		}

		assetURLs, err := findAsset(ctx, c, st, item, pid, siid)
		if err != nil {
			if err == errSkipped {
				log.Infof("searchindex webhook: skipped: all assets are not decompressed or no lod1 bldg")
			} else {
				log.Errorf("searchindex webhook: failed to find asset: %v", err)
			}
			return nil
		}

		if err := c.CommentToItem(ctx, item.ID, "検索インデックスの構築を開始しました。"); err != nil {
			log.Errorf("searchindex webhook: failed to comment: %s", err)
		}

		if _, err := c.UpdateItem(ctx, item.ID, Item{
			SearchIndexStatus: StatusProcessing,
		}.Fields()); err != nil {
			log.Errorf("searchindex webhook: failed to update item: %w", err)
		}

		log.Infof("searchindex webhook: start processing")

		result, err := do(ctx, c, pid, assetURLs)
		if err != nil {
			log.Errorf("searchindex webhook: %v", err)

			if _, err := c.UpdateItem(ctx, item.ID, Item{
				SearchIndexStatus: StatusError,
			}.Fields()); err != nil {
				log.Errorf("searchindex webhook: failed to update item: %s", err)
			}

			if err := c.CommentToItem(ctx, item.ID, fmt.Sprintf("検索インデックスの構築に失敗しました。%v", err)); err != nil {
				log.Errorf("searchindex webhook: failed to comment: %s", err)
			}
			return nil
		}

		if _, err := c.UpdateItem(ctx, item.ID, Item{
			SearchIndexStatus: StatusOK,
			SearchIndex:       result,
		}.Fields()); err != nil {
			log.Errorf("searchindex webhook: failed to update item: %s", err)
		}

		if err := c.CommentToItem(ctx, item.ID, "検索インデックスの構築が完了しました。"); err != nil {
			log.Errorf("searchindex webhook: failed to comment: %s", err)
		}

		log.Infof("searchindex webhook: done")
		return nil
	}, nil
}

func getItem(ctx context.Context, c cms.Interface, st *Storage, w *cmswebhook.Payload) (item Item, siid string, err error) {
	var witem *cms.Item

	if w.Type == cmswebhook.EventAssetDecompress {
		// when asset was decompressed, find data from storage and get the item
		if w.AssetData == nil || w.AssetData.ID == "" {
			log.Debugf("searchindex webhook: invalid event data: %+v", w.Data)
			return
		}

		aid := w.AssetData.ID
		si, err2 := st.FindByAsset(ctx, aid)
		if err2 != nil {
			err = fmt.Errorf("cannot get data from storage: %v", err2)
			return
		} else if si.ID == "" {
			err = errors.New("item and asset not registered to storage")
			return
		}

		siid = si.ID
		witem, err2 = c.GetItem(ctx, si.Item)
		if err2 != nil {
			err = fmt.Errorf("cannot get item %s: %v", si.Item, err2)
			return
		}

		si = si.RemoveAsset(aid)
		if err := st.Set(ctx, si); err != nil {
			log.Errorf("searchindex webook: cannot set to storage: %w", err)
		}
	} else {
		// when item was created or updated
		if w.ItemData == nil || w.ItemData.Item == nil || w.ItemData.Model == nil {
			log.Debugf("searchindex webhook: invalid event data: %+v", w.Data)
			return
		}

		if w.ItemData.Model.Key != modelKey {
			log.Debugf("searchindex webhook: invalid model id: %s, key: %s", w.ItemData.Item.ModelID, w.ItemData.Model.Key)
			return
		}

		witem = w.ItemData.Item
	}

	if witem == nil {
		return
	}

	item = ItemFrom(*witem)
	return
}

func findAsset(ctx context.Context, c cms.Interface, st *Storage, item Item, pid, siid string) ([]*url.URL, error) {
	var assetNotDecompressed []string
	var urls []*url.URL
	for _, aid := range item.Bldg {
		a, err := c.Asset(ctx, aid)
		if err != nil {
			return nil, fmt.Errorf("failed to get an asset (%s): %s", aid, err)
		}

		u, _ := url.Parse(a.URL)
		if u == nil || path.Ext(u.Path) != ".zip" {
			continue
		}

		name := pathFileName(u.Path)
		if !strings.Contains(name, "_lod1") {
			continue
		}

		if a.ArchiveExtractionStatus != cms.AssetArchiveExtractionStatusDone {
			// register asset ID and item ID to storage
			assetNotDecompressed = append(assetNotDecompressed, aid)
			continue
		}

		urls = append(urls, u)
	}

	if len(assetNotDecompressed) > 0 {
		if err := st.Set(ctx, StorageItem{
			ID:    siid,
			Item:  item.ID,
			Asset: assetNotDecompressed,
		}); err != nil {
			return nil, fmt.Errorf("failed to save to storage: %s", err)
		}

		return nil, errSkipped
	}

	if len(urls) == 0 {
		return nil, errSkipped
	}

	return urls, nil
}

func do(ctx context.Context, c cms.Interface, pid string, u []*url.URL) ([]string, error) {
	var results []string
	for _, u := range u {
		name := pathFileName(u.Path)
		if name == "" {
			continue
		}

		log.Infof("searchindex webhook: start processing for %s", name)
		indexer := NewZipIndexer(c, pid, u)
		aid, err := indexer.BuildIndex(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("「%s」の処理中にエラーが発生しました。%w", name, err)
		}
		results = append(results, aid)
	}
	return results, nil
}

func pathFileName(p string) string {
	return strings.TrimSuffix(path.Base(p), path.Ext(p))
}
