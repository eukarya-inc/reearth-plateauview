package preparegspatialjp

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/reearth/reearthx/log"
)

func DownloadFile(ctx context.Context, url string) (*bytes.Reader, error) {
	log.Infofc(ctx, "downloading %s...", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	b := &bytes.Buffer{}
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b.Bytes()), nil
}

func Unzip(ctx context.Context, zipFile *bytes.Reader, targetDir string, trimPathSuffix string) error {
	// Open the zip archive for reading.
	r, err := zip.NewReader(zipFile, int64(zipFile.Len()))
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}

	// Iterate through the files in the archive.
	for _, f := range r.File {
		// Determine the file path ensuring it's within targetDir.
		filePath := filepath.Join(targetDir, strings.TrimPrefix(f.Name, trimPathSuffix))
		filePath = strings.TrimSuffix(filePath, `\`)
		// filePath = filepath.ToSlash(filePath)
		filePath = strings.ReplaceAll(filePath, `\`, `/`)
		log.Infofc(ctx, "unzipping %s -> %s", f.Name, filePath)

		if trimPathSuffix != "" {
			filePath = strings.Trim(filePath, trimPathSuffix)
		}

		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		// If it's a directory, create it.
		if f.FileInfo().IsDir() {
			log.Infofc(ctx, "creating dir %s...", f.Name)
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
			_ = destFile.Close()
			return fmt.Errorf("failed to open file in zip: %v", err)
		}

		// Copy the file contents from the zip archive to the new file.
		if _, err := io.Copy(destFile, srcFile); err != nil {
			_ = destFile.Close()
			_ = srcFile.Close()
			log.Errorfc(ctx, "failed to copy file contents: %v", err)
			continue
			// return fmt.Errorf("failed to copy file contents: %v", err)
		}

		// Close the file descriptors.
		_ = destFile.Close()
		_ = srcFile.Close()
	}
	return nil
}

func ZipDir(ctx context.Context, srcDir string, destZip string) error {
	log.Infofc(ctx, "start zipping %s to %s...", srcDir, destZip)

	file, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		log.Infofc(ctx, "zipping %s...", path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		path = strings.TrimPrefix(path, srcDir)

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, file); err != nil {
			return err
		}

		return nil
	}

	if err := filepath.Walk(srcDir, walker); err != nil {
		return fmt.Errorf("failed to walk dir: %v", err)
	}

	return nil
}
