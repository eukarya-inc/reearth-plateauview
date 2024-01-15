package preparegspatialjp

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cms "github.com/reearth/reearth-cms-api/go"
)

func PrepareCityGML(ctx context.Context, cms *cms.CMS, cityItem *CityItem, allFeatureItems map[string]FeatureItem) (string, string, error) {
	return "", "", nil
	// zipFileName := "citygml.zip"
	tmpPath, outPath, err := prepareTmpDir()
	if err != nil {
		return "", "", err
	}

	// panic("not implemented")
	for _, ft := range featureTypes {
		if fi, ok := allFeatureItems[ft]; ok {
			if fi.CityGML == "" {
				continue
			}

			log.Printf("downloading citygml for %s...", ft)

			downloadPath := tmpPath + ft + ".zip"
			if err := downloadFile(downloadPath, fi.CityGML); err != nil {
				log.Printf("failed to download citygml for %s: %s", ft, err)
				return "", "", err
			}

			if err := unzip(downloadPath, tmpPath); err != nil {
				log.Printf("failed to unzip citygml for %s: %s", ft, err)
				return "", "", err
			}

			if err := moveDir(tmpPath+ft, outPath+ft); err != nil {
				log.Printf("failed to move citygml for %s: %s", ft, err)
				return "", "", err
			}

			os.Remove(downloadPath)

			log.Printf("getting citygml for %s...", ft)
		}
	}

	return "", "", nil
}

func prepareTmpDir() (string, string, error) {
	tmpPath := "tmp/"
	outPath := tmpPath + "citygml/"

	fileInfo, err := os.Lstat("./")

	if err != nil {
		return "", "", err
	}

	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm

	if err := os.MkdirAll(outPath, unixPerms); err != nil {
		return "", "", err
	}

	return tmpPath, outPath, nil
}

func downloadFile(filePath string, url string) error {
	log.Printf("downloading %s...", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func unzip(zipFile, targetDir string) error {
	// Open the zip archive for reading.
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()
	// Iterate through the files in the archive.
	for _, f := range r.File {
		// Determine the file path ensuring it's within targetDir.
		filePath := filepath.Join(targetDir, f.Name)
		filePath = strings.TrimSuffix(filePath, `\`)
		filePath = strings.ReplaceAll(filePath, `\`, `/`)
		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}
		// If it's a directory, create it.
		if f.FileInfo().IsDir() {
			log.Printf("creating dir %s...", f.Name)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create dir: %v", err)
			}
			continue
		}
		// Create the enclosing directory if needed.
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create dir: %v", err)
		}
		// Create the file.
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		// Open the file within the zip archive for reading.
		srcFile, err := f.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open file in zip: %v", err)
		}
		// Copy the file contents from the zip archive to the new file.
		if _, err := io.Copy(destFile, srcFile); err != nil {
			destFile.Close()
			srcFile.Close()
			return fmt.Errorf("failed to copy file contents: %v", err)
		}
		// Close the file descriptors.
		destFile.Close()
		srcFile.Close()
	}
	return nil
}

func moveDir(src string, dest string) error {
	log.Printf("moving %s to %s...", src, dest)

	if err := os.Rename(src, dest); err != nil {
		return err
	}

	return nil
}
