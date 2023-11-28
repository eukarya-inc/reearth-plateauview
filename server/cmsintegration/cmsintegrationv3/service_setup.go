package cmsintegrationv3

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/reearth/reearthx/log"
	"golang.org/x/exp/slices"
)

type SetupCityItemsInput struct {
	ProjectID string `json:"projectId"`
	DataURL   string `json:"dataUrl"`
	Test      bool   `json:"test"`
	Force     bool   `json:"force"`
	Offset    int    `json:"offset"`
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
	if len(modelIDs) == 0 || modelIDs[cityModel] == "" || modelIDs[relatedModel] == "" {
		return fmt.Errorf("no models found")
	}

	cityModel := modelIDs[cityModel]
	relatedModel := modelIDs[relatedModel]

	// check city item total count
	if !inp.Force {
		items, err := s.CMS.GetItemsPartially(ctx, cityModel, 0, 1, false)
		if err != nil {
			return fmt.Errorf("failed to get city items: %w", err)
		}
		if items.TotalCount > 0 {
			return fmt.Errorf("city items already exist")
		}
	}

	// parse data
	setupItems, features, err := getAndParseSetupCSV(ctx, s, inp.DataURL)
	if err != nil {
		return fmt.Errorf("failed to get and parse data: %w", err)
	}

	if inp.Offset > 0 {
		setupItems = setupItems[inp.Offset:]
	}

	// if test is true, setupItems count is limited to 10
	if inp.Test && len(setupItems) > 10 {
		setupItems = setupItems[:10]
	}

	for _, f := range features {
		if modelIDs[f] == "" {
			return fmt.Errorf("model id for %s is not found", f)
		}
	}

	log.Infofc(ctx, "cmsintegrationv3: setup %d items", len(setupItems))

	// process cities
	for i, item := range setupItems {
		if onprogress != nil {
			onprogress(i, len(setupItems))
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
			return fmt.Errorf("failed to create city item (%d/%d): %w", i, len(setupItems), err)
		}

		relatedItem := &RelatedItem{
			City: newCityItem.ID,
		}

		newRelatedItem, err := s.CMS.CreateItem(ctx, relatedModel, relatedItem.CMSItem().Fields, nil)
		if err != nil {
			return fmt.Errorf("failed to create related data item (%d/%d): %w", i, len(setupItems), err)
		}

		featureItemIDs := map[string]string{}
		for _, f := range features {
			var status ManagementStatus
			if !slices.Contains(item.Features, f) {
				status = ManagementStatusSkip
			}

			featureItem := &FeatureItem{
				City:   newCityItem.ID,
				Status: status,
			}
			featureCMSItem := featureItem.CMSItem()

			newFeatureItem, err := s.CMS.CreateItem(ctx, modelIDs[f], featureCMSItem.Fields, featureCMSItem.MetadataFields)
			if err != nil {
				return fmt.Errorf("failed to create feature item (%d/%d/%s): %w", i, len(setupItems), f, err)
			}

			featureItemIDs[f] = newFeatureItem.ID
		}

		if _, err := s.CMS.UpdateItem(ctx, newCityItem.ID, (&CityItem{
			References:     featureItemIDs,
			RelatedDataset: newRelatedItem.ID,
		}).CMSItem().Fields, nil); err != nil {
			return fmt.Errorf("failed to update city item (%d/%d): %w", i, len(setupItems), err)
		}
	}

	return nil
}

func getAndParseSetupCSV(ctx context.Context, s *Services, url string) ([]SetupCSVItem, []string, error) {
	r, err := s.GET(ctx, url)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	return parseSetupCSV(ctx, r)
}

func parseSetupCSV(ctx context.Context, r io.Reader) ([]SetupCSVItem, []string, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// read header
	header, err := cr.Read()
	if err != nil {
		return nil, nil, err
	}
	if len(header) < columnFeaturesIndex+1 {
		return nil, nil, fmt.Errorf("invalid header: %v", header)
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
			return nil, nil, err
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

	return items, features, nil
}
