package preparegspatialjp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func PrepareCityGML(ctx context.Context, cms *cms.CMS, tmpDir string, cityItem *CityItem, allFeatureItems map[string]FeatureItem, uc int) (string, string, error) {

	dataName := fmt.Sprintf("%s_%s_city_%d_citygml_%d_op", cityItem.CityCode, cityItem.CityNameEn, cityItem.YearInt(), uc)
	downloadPath := filepath.Join(tmpDir, dataName)
	_ = os.MkdirAll(downloadPath, os.ModePerm)

	zipFileName := dataName + ".zip"
	zipFilePath := filepath.Join(tmpDir, zipFileName)

	log.Infofc(ctx, "preparing citygml: %s", dataName)
	if err := getAssets(ctx, cms, cityItem, downloadPath, tmpDir); err != nil {
		return "", "", fmt.Errorf("failed to get assets: %w", err)
	}

	if err := getUdx(ctx, allFeatureItems, downloadPath, tmpDir); err != nil {
		return "", "", fmt.Errorf("failed to get udx: %w", err)
	}

	if err := ZipDir(ctx, downloadPath, zipFilePath, false); err != nil {
		return "", "", fmt.Errorf("failed to zip citygml: %w", err)
	}

	return zipFileName, zipFilePath, nil
}

func getUdx(ctx context.Context, allFeatureItems map[string]FeatureItem, dest, tmpDir string) error {
	outPath := filepath.Join(dest, "udx")
	_ = os.MkdirAll(outPath, os.ModePerm)

	for _, ft := range featureTypes {
		fi, ok := allFeatureItems[ft]
		if !ok || fi.CityGML == "" {
			continue
		}

		log.Infofc(ctx, "downloading citygml for %s...", ft)

		if _, err := downloadAndUnzip(ctx, fi.CityGML, outPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip(ft, "udx/"),
		}); err != nil {
			return fmt.Errorf("failed to unzip citygml for %s: %w", ft, err)
		}
	}
	return nil
}

func getAssets(ctx context.Context, cms *cms.CMS, cityItem *CityItem, downloadPath, tmpDir string) error {
	codeLists := cityItem.CodeLists
	if codeLists != "" {
		log.Infofc(ctx, "downloading codeLists: %s...", codeLists)

		assets, err := cms.Asset(ctx, codeLists)
		if err != nil {
			return fmt.Errorf("failed to get assets codeLists: %w", err)
		}

		if _, err := downloadAndUnzip(ctx, assets.URL, downloadPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip("codelists", ""),
		}); err != nil {
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

		if _, err := downloadAndUnzip(ctx, assets.URL, downloadPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip("schemas", ""),
		}); err != nil {
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

		if _, err := downloadAndUnzip(ctx, assets.URL, downloadPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip("metadata", ""),
		}); err != nil {
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

		if _, err := downloadAndUnzip(ctx, assets.URL, downloadPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip("specification", ""),
		}); err != nil {
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

		if _, err := downloadAndUnzip(ctx, assets.URL, downloadPath, tmpDir, &UnzipOptions{
			Rename: renameCityGMLZip("", "misc/"),
		}); err != nil {
			return fmt.Errorf("failed to unzip assets misc: %w", err)
		}
	}

	return nil
}

func renameCityGMLZip(ty, prefix string) func(p string) (string, error) {
	return func(rawPath string) (string, error) {
		p := rawPath
		if prefix != "" {
			if strings.HasPrefix(p, prefix) {
				p = strings.TrimPrefix(p, prefix)
			} else if strings.HasSuffix(prefix, "/") && rawPath == prefix[:len(prefix)-1] {
				return "", SkipUnzip
			}
		}

		if ty != "" {
			paths := strings.Split(p, "/")

			if len(paths) > 0 {
				if strings.HasSuffix(paths[0], "_"+ty) {
					paths[0] = ty
				} else if paths[0] != ty {
					return "", fmt.Errorf("unexpected path: %s", p)
				}
				if len(paths) > 1 && paths[1] == ty {
					// remove paths[1]
					paths = append(paths[:1], paths[2:]...)
				}
			} else {
				return "", fmt.Errorf("unexpected path: %s", p)
			}

			res := strings.Join(paths, "/")
			if res == rawPath {
				return "", nil
			}

			return res, nil
		}

		if p == rawPath {
			return "", nil
		}
		return p, nil
	}
}
