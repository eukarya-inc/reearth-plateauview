package preparegspatialjp

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

var DirMap = map[string]string{
	"codelists":     "コードリスト",
	"metadata":      "メタデータ",
	"schemas":       "CityGMLスキーマ",
	"specification": "東京23区における3D都市モデルのための拡張製品仕様書",
	"indexmap":      "索引図 (PDF)",
	"bldg":          "建築物 (CityGML)",
	"tran":          "道路 (CityGML)",
	"urf":           "都市計画決定情報 (CityGML)",
	"luse":          "土地利用 (CityGML)",
	"dem":           "地形 (CityGML)",
	"frn":           "都市設備 (CityGML)",
	"bird":          "汎用オブジェクト (CityGML)",
	"lsld":          "土砂災害警戒区域 (CityGML)",
	"htd":           "高潮浸水想定区域 (CityGML)",
	"fld":           "洪水浸水想定区域 (CityGML)",
	"natl":          "国管理 (CityGML)",
	"pref":          "都道府県管理 (CityGML)",
}

func GenerateCityGMLMarkdown(ctx context.Context, zipFileName, zipPath string) (string, error) {
	indexPath := filepath.Join("tmp", "citygml.md")

	destFile, err := os.OpenFile(indexPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	md, err := ZipToMarkdownTree(ctx, zipFileName, zipPath)
	if err != nil {
		return "", err
	}

	_, err = destFile.WriteString(md)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return md, nil
}

func ZipToMarkdownTree(ctx context.Context, zipFileName, zipPath string) (string, error) {
	log.Infofc(ctx, "start generating markdown %s...", zipPath)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	structure := map[string][]string{}

	// Iterate through each file in the zip archive
	for _, f := range r.File {
		dir, file := filepath.Split(f.Name)
		structure[dir] = append(structure[dir], file)
	}

	{
		pp := pp.New()
		pp.SetColoringEnabled(false)
		s := pp.Sprint(structure)
		log.Infofc(ctx, "structure: %s", s)
	}

	zipFileSize, err := getFileSizeKBMBGB(zipPath)
	if err != nil {
		return "", err
	}

	// Use a string builder to construct the Markdown output
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("- %s：CityGML（v2）(%s)\n", zipFileName, zipFileSize))

	var finishedDirs []string

	subDirs := lo.Keys(structure)
	sort.Strings(subDirs)

	for _, subDir := range subDirs {
		fileNames := structure[subDir]
		splitSubDir := strings.Split(subDir, "/")

		log.Infofc(ctx, "subDir: %s", subDir)
		log.Infofc(ctx, "finishedDirs: %+v", finishedDirs)
		if finished := lo.SomeBy(finishedDirs, func(d string) bool {
			return strings.HasPrefix(subDir, d)
		}); finished {
			log.Infofc(ctx, "finished: %s", subDir)
			continue
		}

		index1 := strings.Repeat("  ", 1)
		index2 := strings.Repeat("  ", 2)
		index3 := strings.Repeat("  ", 3)

		if subDir == "/" {
			const indexmap = "indexmap"
			ok := lo.SomeBy(fileNames, func(x string) bool {
				return strings.Contains(x, indexmap)
			})

			if ok {
				sb.WriteString(fmt.Sprintf("%s- %s：索引図 (PDF)\n", index1, indexmap))
				continue
			}
		}

		var index1Dirs = []string{"/codelists/", "/metadata/", "/schemas/", "/specification/"}

		if ok := lo.SomeBy(index1Dirs, func(d string) bool {
			return strings.HasPrefix(subDir, d)
		}); ok {
			log.Infofc(ctx, "index1Dirs: %s", subDir)

			dir := splitSubDir[1]
			explain, ok := DirMap[dir]
			if !ok {
				log.Warnfc(ctx, "index1Dirs not found: %s", dir)
				continue
			}
			sb.WriteString(fmt.Sprintf("%s- %s：%s\n", index1, dir, explain))
			finishedDirs = append(finishedDirs, "/"+dir)
			continue
		}

		if ok := strings.HasPrefix(subDir, "/udx/"); ok {

			if firstUDX := !lo.SomeBy(finishedDirs, func(d string) bool {
				return strings.HasPrefix(d, "/udx/")
			}); firstUDX {
				log.Infofc(ctx, "udx first: %s", subDir)
				sb.WriteString(fmt.Sprintf("%s- udx：\n", index1))
			}

			ft := splitSubDir[2]
			udxFt := "/udx/" + ft

			if ft == "fld" {
				if firstFLD := !lo.SomeBy(finishedDirs, func(d string) bool {
					return strings.HasPrefix(d, "/udx/fld/")
				}); firstFLD {
					log.Infofc(ctx, "fld first: %s", subDir)
					sb.WriteString(fmt.Sprintf("%s- fld：\n", index2))
				}
				subFld := splitSubDir[3]
				explain, ok := DirMap[subFld]
				if !ok {
					log.Warnfc(ctx, "fld not found: %s", ft)
					continue
				}

				first := !lo.SomeBy(finishedDirs, func(d string) bool {
					return strings.HasPrefix(subDir, d)
				})
				if !first {
					log.Infofc(ctx, "udx already exists: %s", ft)
					continue
				}

				sb.WriteString(fmt.Sprintf("%s- %s：%s\n", index3, subFld, explain))
				finishedDirs = append(finishedDirs, udxFt+"/"+subFld)
				continue
			}

			explain, ok := DirMap[ft]
			if !ok {
				log.Warnfc(ctx, "udx not found: %s", ft)
				continue
			}

			first := !lo.SomeBy(finishedDirs, func(d string) bool {
				return strings.HasPrefix(subDir, d)
			})
			if !first {
				log.Infofc(ctx, "udx already exists: %s", ft)
				continue
			}

			sb.WriteString(fmt.Sprintf("%s- %s：%s\n", index2, ft, explain))
			finishedDirs = append(finishedDirs, udxFt)
		}
	}

	return sb.String(), nil
}

func getFileSizeKBMBGB(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSizeBytes := fileInfo.Size()

	// Convert size into MB or GB
	const (
		_  = iota // ignore first value by assigning to blank identifier
		KB = 1 << (10 * iota)
		MB
		GB
	)

	var sizeStr string
	switch {
	case fileSizeBytes >= GB:
		sizeStr = fmt.Sprintf("%.2f GB", float64(fileSizeBytes)/GB)
	case fileSizeBytes >= MB:
		sizeStr = fmt.Sprintf("%.2f MB", float64(fileSizeBytes)/MB)
	case fileSizeBytes >= KB:
		sizeStr = fmt.Sprintf("%.2f KB", float64(fileSizeBytes)/KB)
	default:
		sizeStr = fmt.Sprintf("%d Bytes", fileSizeBytes)
	}

	return sizeStr, nil
}
