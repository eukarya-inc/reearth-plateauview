package geospatialjpv3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/vincent-petithory/dataurl"
)

type Seed struct {
	CityGML            string
	Plateau            string
	Related            string
	Desc               string
	CityGMLDescription string
	PlateauDescription string
	RelatedDescription string
	Area               string
	ThumbnailURL       string `pp:"-"`
	Version            int
	Author             string
	AuthorEmail        string
	Maintainer         string
	MaintainerEmail    string
	Quality            string
}

const defaultVersion = 3

func getSeed(ctx context.Context, c cms.Interface, cityItem *CityItem) (Seed, error) {
	seed := Seed{
		Version: defaultVersion,
	}

	rawDataItem, err := c.GetItem(ctx, cityItem.GeospatialjpData, true)
	if err != nil {
		return seed, fmt.Errorf("failed to get data item: %w", err)
	}

	rawIndexItem, err := c.GetItem(ctx, cityItem.GeospatialjpIndex, true)
	if err != nil {
		return seed, fmt.Errorf("failed to get index item: %w", err)
	}

	var dataItem CMSDataItem
	rawDataItem.Unmarshal(&dataItem)

	var indexItem CMSIndexItem
	rawIndexItem.Unmarshal(&indexItem)

	log.Debugfc(ctx, "geospatialjpv3: rawDataItem: %s", ppp.Sprint(rawDataItem))
	log.Debugfc(ctx, "geospatialjpv3: rawIndexItem: %s", ppp.Sprint(rawIndexItem))
	log.Debugfc(ctx, "geospatialjpv3: dataItem: %s", ppp.Sprint(dataItem))
	log.Debugfc(ctx, "geospatialjpv3: indexItem: %s", ppp.Sprint(indexItem))

	if dataItem.CityGML != nil {
		seed.CityGML = valueToAsset(dataItem.CityGML)
	}
	if dataItem.Plateau != nil {
		seed.Plateau = valueToAsset(dataItem.Plateau)
	}
	if dataItem.Related != nil {
		seed.Related = valueToAsset(dataItem.Related)
	}
	if indexItem.Desc != "" {
		seed.Desc = indexItem.Desc
	}

	seed.CityGMLDescription = indexItem.DescCityGML
	seed.PlateauDescription = indexItem.DescPlateau
	seed.RelatedDescription = indexItem.DescRelated
	seed.Area = indexItem.Region

	if thumnailURL := valueToAsset(indexItem.Thumbnail); thumnailURL != "" {
		seed.ThumbnailURL, err = fetchAndGetDataURL(thumnailURL)
		if err != nil {
			return seed, fmt.Errorf("failed to fetch thumnail: %w", err)
		}
	}

	return seed, nil
}

func valueToAsset(v map[string]any) string {
	if v == nil {
		return ""
	}
	if url, ok := v["url"].(string); ok {
		return url
	}
	return ""
}

func fetchAndGetDataURL(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch thumnail: %s", res.Status)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", err
	}

	data := buf.Bytes()
	mediaType := http.DetectContentType(data)
	if !strings.HasPrefix(mediaType, "image/") {
		return "", fmt.Errorf("thumnail is not image")
	}

	return dataurl.New(data, mediaType).String(), nil
}
