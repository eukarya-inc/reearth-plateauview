package preparegspatialjp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func PreparePlateau(ctx context.Context, cms *cms.CMS, cityItem *CityItem, allFeatureItems map[string]FeatureItem) (string, string, error) {
	log.Infofc(ctx, "preparing plateau...")

	tmpDir := "tmp"
	downloadPath := filepath.Join(tmpDir, cityItem.CityCode+"_"+cityItem.CityNameEn+"_plateau")
	_ = os.MkdirAll(downloadPath, os.ModePerm)

	zipFileName := fmt.Sprintf("%s_%s_city_%d_3dtiles_mvt.zip", cityItem.CityCode, cityItem.CityNameEn, cityItem.YearInt())
	zipFilePath := filepath.Join(tmpDir, zipFileName)

	for _, ft := range featureTypes {
		fi, ok := allFeatureItems[ft]
		if !ok || fi.Data == nil {
			log.Infofc(ctx, "no data for %s", ft)
			continue
		}

		log.Infofc(ctx, "downloading data for %s...", ft)

		for _, url := range fi.Data {
			log.Infofc(ctx, "downloading url: %s", url)

			if url == "" {
				continue
			}

			data, err := downloadFileAsByteReader(ctx, url)
			if err != nil {
				return "", "", fmt.Errorf("failed to download data for %s: %w", ft, err)
			}

			if err := Unzip(ctx, data, downloadPath, ""); err != nil {
				return "", "", fmt.Errorf("failed to unzip data for %s: %w", ft, err)
			}
		}
	}

	if err := ZipDir(ctx, downloadPath, zipFilePath); err != nil {
		return "", "", fmt.Errorf("failed to zip plateau: %w", err)
	}

	return zipFileName, zipFilePath, nil
}
