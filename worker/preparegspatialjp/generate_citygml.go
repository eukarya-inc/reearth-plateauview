package preparegspatialjp

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/reearth/reearthx/log"
)

var citygmlFiles = []string{
	"codelists",
	"schemas",
	"metadata",
	"specification",
	"misc",
}

func PrepareCityGML(ctx context.Context, tmpDir string, cityItem *CityItem, allFeatureItems map[string]FeatureItem, uc int) (string, string, error) {
	// create a zip file
	rootName := fmt.Sprintf("%s_%s_city_%d_citygml_%d_op", cityItem.CityCode, cityItem.CityNameEn, cityItem.YearInt(), uc)
	downloadPath := filepath.Join(tmpDir, rootName)
	_ = os.MkdirAll(downloadPath, os.ModePerm)

	zipFileName := rootName + ".zip"
	f, err := os.Create(zipFileName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}

	defer f.Close()
	zw := zip.NewWriter(f)
	cz := NewCityGMLZipWriter(zw, rootName)
	defer cz.Close()

	// copy files
	for _, ty := range citygmlFiles {
		url := getCityGMLURL(cityItem, ty)
		if url == "" {
			continue
		}

		log.Infofc(ctx, "preparing citygml (%s): %s", ty, rootName)

		prefix := ""
		if ty == "misc" {
			ty = ""
			prefix = "misc/"
		}

		err := cz.DownloadAndWrite(ctx, url, downloadPath, ty, prefix)
		if err != nil {
			return "", "", fmt.Errorf("failed to download and write %s: %w", ty, err)
		}
	}

	// copy features
	for ty, a := range allFeatureItems {
		url := a.CityGML
		if a.CityGML == "" {
			continue
		}

		log.Infofc(ctx, "preparing citygml (%s): %s", ty, rootName)

		if err := cz.DownloadAndWrite(ctx, url, downloadPath, ty, "udx/"); err != nil {
			return "", "", fmt.Errorf("failed to download and write citygml for %s: %w", ty, err)
		}
	}

	zipFilePath := filepath.Join(tmpDir, zipFileName)
	return zipFileName, zipFilePath, nil
}

func getCityGMLURL(item *CityItem, ty string) string {
	switch ty {
	case "codelists":
		return item.CodeLists
	case "schemas":
		return item.Schemas
	case "metadata":
		return item.Metadata
	case "specification":
		return item.Specification
	case "misc":
		return item.Misc
	}
	return ""
}

type CityGMLZipWriter struct {
	w    *Zip2zip
	name string
}

func NewCityGMLZipWriter(w *zip.Writer, name string) *CityGMLZipWriter {
	return &CityGMLZipWriter{
		w:    NewZip2zip(w),
		name: name,
	}
}

func (z *CityGMLZipWriter) Close() error {
	return z.w.Close()
}

func (z *CityGMLZipWriter) DownloadAndWrite(ctx context.Context, url, tempdir, ty, prefix string) error {
	if url == "" {
		return nil
	}

	err := downloadAndConsumeZip(ctx, url, tempdir, func(zr *zip.Reader, fi os.FileInfo) error {
		log.Debugfc(ctx, "downloaded %s (%s)", url, humanize.Bytes(uint64(fi.Size())))
		reportDiskUsage(tempdir)

		return z.Write(ctx, zr, ty, prefix)
	})

	if err != nil {
		return err
	}

	reportDiskUsage(tempdir)
	return nil
}

func (z *CityGMLZipWriter) Write(ctx context.Context, src *zip.Reader, ty, prefix string) error {
	fn := cityGMLZipPath(ty, prefix)
	return z.w.Run(src, func(f *zip.File) (string, error) {
		p, err := fn(f.Name)
		if err != nil {
			return "", err
		}

		if p == "" {
			log.Debugfc(ctx, "zipping %s: %s -> [SKIP]", z.name, f.Name)
			return "", nil
		}

		log.Debugfc(ctx, "zipping %s: %s -> %s", z.name, f.Name, p)
		return p, nil
	})
}

func cityGMLZipPath(ty, prefix string) func(string) (string, error) {
	return func(rawPath string) (string, error) {
		p := normalizeZipFilePath(rawPath)
		if p == "" {
			return "", nil
		}

		if prefix != "" {
			if strings.HasPrefix(p, prefix) {
				p = strings.TrimPrefix(p, prefix)
			} else if strings.HasSuffix(prefix, "/") && rawPath == prefix[:len(prefix)-1] {
				return "", nil
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
