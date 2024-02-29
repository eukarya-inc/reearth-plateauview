package preparegspatialjp

import (
	"archive/zip"
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

	f, err := os.Create(zipFilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}

	defer f.Close()

	z := NewZip2zip(zip.NewWriter(f))
	defer z.Close()

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

			err := downloadAndConsumeZip(ctx, url, downloadPath, func(zr *zip.Reader, _ os.FileInfo) error {
				return z.Run(zr, func(f *zip.File) (string, error) {
					p := normalizeZipFilePath(f.Name)
					return p, nil
				})
			})
			if err != nil {
				return "", "", fmt.Errorf("failed to download and consume zip: %w", err)
			}
		}
	}

	return zipFileName, zipFilePath, nil
}
