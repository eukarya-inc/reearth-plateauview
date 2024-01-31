package preparegspatialjp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func PrepareCityGML(ctx context.Context, cms *cms.CMS, cityItem *CityItem, allFeatureItems map[string]FeatureItem) (string, string, string, string, error) {
	tmpDir := "tmp"
	downloadPath := filepath.Join(tmpDir, cityItem.CityCode+"_"+cityItem.CityNameEn+"_citygml")
	_ = os.MkdirAll(downloadPath, os.ModePerm)

	zipFileName := cityItem.CityCode + "_" + cityItem.CityNameEn + "_city_2023_citygml_1_op.zip"
	zipFilePath := filepath.Join(tmpDir, zipFileName)

	if err := getAssets(ctx, cms, cityItem, downloadPath); err != nil {
		return "", "", "", "", fmt.Errorf("failed to get assets: %w", err)
	}

	if err := getUdx(ctx, allFeatureItems, downloadPath); err != nil {
		return "", "", "", "", fmt.Errorf("failed to get udx: %w", err)
	}

	if err := ZipDir(ctx, downloadPath, zipFilePath); err != nil {
		return "", "", "", "", fmt.Errorf("failed to zip citygml: %w", err)
	}

	md, err := ZipToMarkdownTree(ctx, zipFileName, zipFilePath)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate markdown: %w", err)
	}

	mdFileName := "citygml.md"
	mdFilePath := filepath.Join(tmpDir, mdFileName)

	mdFile, err := os.OpenFile(mdFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to create file: %v", err)
	}
	_, err = mdFile.WriteString(md)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to write file: %v", err)
	}

	return zipFileName, zipFilePath, mdFileName, mdFilePath, nil
}

func getUdx(ctx context.Context, allFeatureItems map[string]FeatureItem, downloadPath string) error {
	outPath := filepath.Join(downloadPath, "udx")
	_ = os.MkdirAll(outPath, os.ModePerm)

	for _, ft := range featureTypes {
		fi, ok := allFeatureItems[ft]
		if !ok || fi.CityGML == "" {
			continue
		}

		log.Infofc(ctx, "downloading citygml for %s...", ft)

		data, err := DownloadFile(ctx, fi.CityGML)
		if err != nil {
			return fmt.Errorf("failed to download citygml for %s: %w", ft, err)
		}

		if err := Unzip(ctx, data, outPath, ""); err != nil {
			return fmt.Errorf("failed to unzip citygml for %s: %w", ft, err)
		}
	}
	return nil
}

func getAssets(ctx context.Context, cms *cms.CMS, cityItem *CityItem, downloadPath string) error {
	codeLists := cityItem.CodeLists
	if codeLists != "" {
		log.Infofc(ctx, "downloading codeLists: %s...", codeLists)

		assets, err := cms.Asset(ctx, codeLists)
		if err != nil {
			return fmt.Errorf("failed to get assets codeLists: %w", err)
		}

		data, err := DownloadFile(ctx, assets.URL)
		if err != nil {
			return fmt.Errorf("failed to download assets codeLists: %w", err)
		}

		if err := Unzip(ctx, data, downloadPath, ""); err != nil {
			return fmt.Errorf("failed to unzip assets codeLists: %w", err)
		}
	}

	schemas := cityItem.Schemas
	if schemas != "" {
		log.Infofc(ctx, "downloading schemas: %s...", schemas)

		assets, err := cms.Asset(ctx, schemas)
		if err != nil {
			return fmt.Errorf("failed to get assets schemas: %w", err)
		}

		data, err := DownloadFile(ctx, assets.URL)
		if err != nil {
			return fmt.Errorf("failed to download assets schemas: %w", err)
		}

		if err := Unzip(ctx, data, downloadPath, ""); err != nil {
			return fmt.Errorf("failed to unzip assets schemas: %w", err)
		}
	}

	metadata := cityItem.Metadata
	if metadata != "" {
		log.Infofc(ctx, "downloading metadata: %s...", metadata)

		assets, err := cms.Asset(ctx, metadata)
		if err != nil {
			return fmt.Errorf("failed to get assets metadata: %w", err)
		}

		data, err := DownloadFile(ctx, assets.URL)
		if err != nil {
			return fmt.Errorf("failed to download assets metadata: %w", err)
		}

		if err := Unzip(ctx, data, downloadPath, ""); err != nil {
			return fmt.Errorf("failed to unzip assets metadata: %w", err)
		}
	}

	specification := cityItem.Specification
	if specification != "" {
		log.Infofc(ctx, "downloading specification: %s...", specification)

		assets, err := cms.Asset(ctx, specification)
		if err != nil {
			return fmt.Errorf("failed to get assets specification: %w", err)
		}

		data, err := DownloadFile(ctx, assets.URL)
		if err != nil {
			return fmt.Errorf("failed to download assets specification: %w", err)
		}

		if err := Unzip(ctx, data, downloadPath, ""); err != nil {
			return fmt.Errorf("failed to unzip assets specification: %w", err)
		}
	}

	misc := cityItem.Misc
	if misc != "" {
		log.Infofc(ctx, "downloading misc: %s...", misc)

		assets, err := cms.Asset(ctx, misc)
		if err != nil {
			return fmt.Errorf("failed to get assets misc: %w", err)
		}

		data, err := DownloadFile(ctx, assets.URL)
		if err != nil {
			return fmt.Errorf("failed to download assets misc: %w", err)
		}

		if err := Unzip(ctx, data, downloadPath, "misc"); err != nil {
			return fmt.Errorf("failed to unzip assets misc: %w", err)
		}
	}

	return nil
}
