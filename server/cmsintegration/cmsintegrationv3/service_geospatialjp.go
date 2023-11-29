package cmsintegrationv3

import (
	"context"
	"fmt"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

const geospatialjpPrepare = "geospatialjp_prepare"
const geospatialjpPublish = "geospatialjp_publish"

func handleGeospatialjp(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	// if event type is "item.create" and payload is metadata, skip it
	if w.Type == cmswebhook.EventItemCreate && w.ItemData.Item.OriginalItemID != nil ||
		w.ItemData == nil || w.ItemData.Item == nil || w.ItemData.Model == nil ||
		w.ItemData.Item.FieldByKey(relatedConvStatus) == nil {
		return nil
	}

	if w.ItemData.Model.Key != modelPrefix+cityModel {
		log.Debugfc(ctx, "cmsintegrationv3: not city model")
		return nil
	}

	var prepare, publish bool
	if w.Type == cmswebhook.EventItemUpdate {
		if p1, ok1 := FindFieldChangeByKey(w, geospatialjpPrepare); ok1 {
			prepare, _ = p1.(bool)
		}
		if p2, ok2 := FindFieldChangeByKey(w, geospatialjpPrepare); ok2 {
			publish, _ = p2.(bool)
		}
	} else {
		prepare = lo.FromPtr(w.ItemData.Item.FieldByKey(geospatialjpPrepare).GetValue().Bool())
		publish = lo.FromPtr(w.ItemData.Item.FieldByKey(geospatialjpPublish).GetValue().Bool())
	}

	if !prepare && !publish {
		log.Debugfc(ctx, "cmsintegrationv3: skipped")
		return nil
	}

	mainItem, err := s.GetMainItemWithMetadata(ctx, w.ItemData.Item)
	if err != nil {
		return err
	}

	city := CityItemFrom(mainItem)
	if city == nil || city.GeospatialjpData == "" || city.GeospatialjpIndex == "" {
		log.Debugfc(ctx, "cmsintegrationv3: geospatialjp items not linked")
		return nil
	}

	dataItem, err := s.CMS.GetItem(ctx, city.GeospatialjpData, false)
	if err != nil {
		return fmt.Errorf("failed to get geospatialjp data item: %w", err)
	}

	data := GeospatialjpDataItemFrom(dataItem)

	indexItem, err := s.CMS.GetItem(ctx, city.GeospatialjpIndex, false)
	if err != nil {
		return fmt.Errorf("failed to get geospatialjp index item: %w", err)
	}

	index := GeospatialjpIndexItemFrom(indexItem)

	if data == nil || index == nil {
		log.Debugfc(ctx, "cmsintegrationv3: invalid geospatialjp items")
		return nil
	}

	if prepare {
		if err := preparePackagesForGeospatialjp(ctx, s, city, data, index); err != nil {
			return err
		}
	}

	if publish {
		if err := publishPackagesForGeospatialjp(ctx, s, city, data, index); err != nil {
			return err
		}
	}

	return nil
}

func preparePackagesForGeospatialjp(ctx context.Context, s *Services, city *CityItem, data *GeospatialjpDataItem, index *GeospatialjpIndexItem) error {
	log.Infofc(ctx, "cmsintegrationv3: preparePackagesForGeospatialjp")

	return nil
}

func publishPackagesForGeospatialjp(ctx context.Context, s *Services, city *CityItem, data *GeospatialjpDataItem, index *GeospatialjpIndexItem) error {
	log.Infofc(ctx, "cmsintegrationv3: publishPackagesForGeospatialjp")

	return nil
}
