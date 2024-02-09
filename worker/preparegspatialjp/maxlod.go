package preparegspatialjp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func MergeMaxLOD(ctx context.Context, cms *cms.CMS, tmpDir string, cityItem *CityItem, allFeatureItems map[string]FeatureItem) (string, string, error) {
	log.Infofc(ctx, "preparing plateau...")

	_ = os.MkdirAll(tmpDir, os.ModePerm)

	fileName := fmt.Sprintf("%s_%s_%d_maxlod.csv", cityItem.CityCode, cityItem.CityNameEn, cityItem.YearInt())
	filePath := filepath.Join(tmpDir, fileName)

	allData := bytes.NewBuffer(nil)

	first := false
	for _, ft := range featureTypes {
		fi, ok := allFeatureItems[ft]
		if !ok || fi.MaxLOD == "" {
			log.Infofc(ctx, "no maxlod for %s", ft)
			continue
		}

		log.Infofc(ctx, "downloading maxlod data for %s: %s", ft, fi.MaxLOD)
		data, err := downloadFile(ctx, fi.MaxLOD)
		if err != nil {
			return "", "", fmt.Errorf("failed to download data for %s: %w", ft, err)
		}

		b := bufio.NewReader(data)
		if first {
			if line, err := b.ReadString('\n'); err != nil { // skip the first line
				return "", "", fmt.Errorf("failed to read first line: %w", err)
			} else if line == "" || isNumeric(rune(line[0])) {
				// the first line shold be header (code,type,maxlod,filename)
				return "", "", fmt.Errorf("invalid maxlod data for %s", ft)
			}
		} else {
			first = true
		}

		if _, err := allData.ReadFrom(b); err != nil {
			return "", "", fmt.Errorf("failed to read data for %s: %w", ft, err)
		}
	}

	if allData.Len() > 0 {
		if err := os.WriteFile(filePath, allData.Bytes(), os.ModePerm); err != nil {
			return "", "", fmt.Errorf("failed to write data to file: %w", err)
		}
	} else {
		return "", "", nil
	}

	return fileName, filePath, nil
}

func isNumeric(s rune) bool {
	return strings.ContainsRune("0123456789", s)
}
