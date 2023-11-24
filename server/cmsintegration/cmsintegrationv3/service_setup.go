package cmsintegrationv3

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/reearth/reearthx/log"
)

type SetupCityItemsInput struct {
	ProjectID string `json:"projectId"`
	DataURL   string `json:"dataUrl"`
}

type SetupCSVItem struct {
	Prefecture string   `json:"prefecture"`
	Name       string   `json:"name"`
	NameEn     string   `json:"nameEn"`
	Code       string   `json:"code"`
	Features   []string `json:"features"`
}

const columnFeaturesIndex = 4

func SetupCityItems(ctx context.Context, s *Services, inp SetupCityItemsInput, onprogress func(i, l int)) error {
	if inp.ProjectID == "" {
		return fmt.Errorf("modelId is required")
	}

	if inp.DataURL == "" {
		return fmt.Errorf("dataUrl is required")
	}

	log.Infofc(ctx, "cmsintegrationv3: setup city items to %s", inp.ProjectID)

	// get model info
	modelIDs := map[string]string{}
	models, err := s.CMS.GetModels(ctx, inp.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to get models: %w", err)
	}
	for _, m := range models.Models {
		if !strings.HasPrefix(m.Key, modelPrefix) {
			continue
		}
		modelIDs[strings.TrimPrefix(m.Key, modelPrefix)] = m.ID
	}
	if len(modelIDs) == 0 || modelIDs[modelPrefix+cityModel] == "" || modelIDs[modelPrefix+relatedModel] == "" {
		return fmt.Errorf("no models found")
	}

	cityModel := modelIDs[modelPrefix+cityModel]
	relatedModel := modelIDs[modelPrefix+relatedModel]

	// parse data
	items, err := getAndParseSetupCSV(ctx, s, inp.DataURL)
	if err != nil {
		return fmt.Errorf("failed to get and parse data: %w", err)
	}

	log.Infofc(ctx, "cmsintegrationv3: setup %d items", len(items))

	// process cities
	for i, item := range items {
		if onprogress != nil {
			onprogress(i, len(items))
		}

		cityItem := &CityItem{
			Prefecture: item.Prefecture,
			CityName:   item.Name,
			CityNameEn: item.NameEn,
			CityCode:   item.Code,
		}
		cityCMSItem := cityItem.CMSItem()

		newCityItem, err := s.CMS.CreateItem(ctx, cityModel, cityCMSItem.Fields, cityCMSItem.MetadataFields)
		if err != nil {
			return fmt.Errorf("failed to create city item (%d/%d): %w", i, len(items), err)
		}

		relatedItem := &RelatedItem{
			City: newCityItem.ID,
		}

		newRelatedItem, err := s.CMS.CreateItem(ctx, relatedModel, relatedItem.CMSItem().Fields, nil)
		if err != nil {
			return fmt.Errorf("failed to create related data item (%d/%d): %w", i, len(items), err)
		}

		featureItemIDs := map[string]string{}
		for _, f := range item.Features {
			if modelIDs[f] == "" {
				return fmt.Errorf("model id for %s is not found (%d/%d)", f, i, len(items))
			}

			featureItem := &FeatureItem{
				City:   newCityItem.ID,
				Status: ManagementStatusNotStarted,
			}
			featureCMSItem := featureItem.CMSItem()

			newFeatureItem, err := s.CMS.CreateItem(ctx, modelIDs[f], featureCMSItem.Fields, featureCMSItem.MetadataFields)
			if err != nil {
				return fmt.Errorf("failed to create feature item (%d/%d/%s): %w", i, len(items), f, err)
			}

			featureItemIDs[f] = newFeatureItem.ID
		}

		if _, err := s.CMS.UpdateItem(ctx, newCityItem.ID, (&CityItem{
			References:     featureItemIDs,
			RelatedDataset: newRelatedItem.ID,
		}).CMSItem().Fields, nil); err != nil {
			return fmt.Errorf("failed to update city item (%d/%d): %w", i, len(items), err)
		}
	}

	return nil
}

func getAndParseSetupCSV(ctx context.Context, s *Services, url string) ([]SetupCSVItem, error) {
	r, err := s.GET(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return parseSetupCSV(ctx, r)
}

func parseSetupCSV(ctx context.Context, r io.Reader) ([]SetupCSVItem, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// read header
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	if len(header) < columnFeaturesIndex+1 {
		return nil, fmt.Errorf("invalid header: %v", header)
	}

	features := make([]string, 0, len(header)-3)
	for i := columnFeaturesIndex; i < len(header); i++ {
		features = append(features, header[i])
	}

	var items []SetupCSVItem
	i := 0
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		itemFeatures := make([]string, 0, len(row)-3)
		for i := columnFeaturesIndex; i < len(row); i++ {
			if row[i] != "" {
				itemFeatures = append(itemFeatures, features[i-columnFeaturesIndex])
			}
		}

		items = append(items, SetupCSVItem{
			Name:       row[0],
			NameEn:     row[1],
			Code:       row[2],
			Prefecture: row[3],
			Features:   itemFeatures,
		})
		i++
	}

	return items, nil
}
