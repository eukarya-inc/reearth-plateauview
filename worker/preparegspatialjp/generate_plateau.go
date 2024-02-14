package preparegspatialjp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func PreparePlateau(ctx context.Context, cms *cms.CMS, tmpDir string, cityItem *CityItem, allFeatureItems map[string]FeatureItem, uc int) (string, string, error) {

	dataName := fmt.Sprintf("%s_%s_city_%d_3dtiles_mvt_%d_op", cityItem.CityCode, cityItem.CityNameEn, cityItem.YearInt(), uc)
	downloadPath := filepath.Join(tmpDir, dataName)
	_ = os.MkdirAll(downloadPath, os.ModePerm)

	zipFileName := dataName + ".zip"
	zipFilePath := filepath.Join(tmpDir, zipFileName)

	log.Infofc(ctx, "preparing plateau: %s", dataName)
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

			if _, err := downloadAndUnzip(ctx, url, downloadPath, tmpDir, nil); err != nil {
				return "", "", fmt.Errorf("failed to unzip data for %s: %w", ft, err)
			}
		}
	}

	if err := ZipDir(ctx, downloadPath, zipFilePath, false); err != nil {
		return "", "", fmt.Errorf("failed to zip plateau: %w", err)
	}

	return zipFileName, zipFilePath, nil
}
