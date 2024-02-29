package preparegspatialjp

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/reearth/reearthx/log"
)

func PreparePlateau(ctx context.Context, c *CMSWrapper, m MergeContext) (res string, err error) {
	defer func() {
		err = fmt.Errorf("3D Tiles,MVTのマージに失敗しました: %w", err)
		c.NotifyError(ctx, err, false, true, false)
	}()

	_, path, err := mergePlateau(ctx, m)
	if err != nil {
		err = fmt.Errorf("failed to prepare plateau: %w", err)
		return
	}

	aid, err := c.UploadFile(ctx, path)
	if err != nil {
		err = fmt.Errorf("failed to upload file: %w", err)
		return
	}

	if err2 := c.UpdateDataItem(ctx, &GspatialjpDataItem{
		MergePlateauStatus: successTag,
		Plateau:            aid,
	}); err2 != nil {
		err = fmt.Errorf("failed to update data item: %w", err2)
		return
	}

	res = path
	return
}

func mergePlateau(ctx context.Context, m MergeContext) (string, string, error) {
	tmpDir := m.TmpDir
	cityItem := m.CityItem
	allFeatureItems := m.AllFeatureItems
	uc := m.UC

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

	cz := NewZip2zip(zip.NewWriter(f))
	defer cz.Close()

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
				return cz.Run(zr, func(f *zip.File) (string, error) {
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
