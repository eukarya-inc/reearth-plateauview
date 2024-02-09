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

func downloadFileAsByteReader(ctx context.Context, url string) (*bytes.Reader, error) {
	b, err := downloadFileAsBytes(ctx, url)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func downloadFileAsBytes(ctx context.Context, url string) ([]byte, error) {
	r, err := downloadFile(ctx, url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()

	b := &bytes.Buffer{}
	_, err = io.Copy(b, r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func downloadFile(ctx context.Context, url string) (io.ReadCloser, error) {
	log.Infofc(ctx, "downloading %s...", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	return resp.Body, nil
	// b := &bytes.Buffer{}
	// _, err = io.Copy(b, resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// return bytes.NewReader(b.Bytes()), nil
}

func Unzip(ctx context.Context, zipFile *bytes.Reader, targetDir, prefix string, checkFile func(string) error) error {
	// Open the zip archive for reading.
	r, err := zip.NewReader(zipFile, int64(zipFile.Len()))
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}

	// Iterate through the files in the archive.
	for _, f := range r.File {
		filename := f.Name
		filename = strings.ReplaceAll(filename, `\`, "/")
		filename = strings.TrimSuffix(filename, "/")

		if prefix != "" {
			if strings.HasPrefix(filename, prefix) {
				filename = strings.TrimPrefix(filename, prefix)
			} else {
				log.Debugf("skip %s (no prefix: %s)", f.Name, prefix)
				continue
			}
		}

		if strings.HasPrefix(filename, "__MACOSX/") ||
			strings.HasSuffix(filename, "/.DS_Store") ||
			strings.HasSuffix(filename, "/Thumb.db") ||
			filename == ".DS_Store" || filename == "Thumbs.db" {
			log.Debugf("skip %s (%s)", f.Name, filename)
			continue
		}

		if checkFile != nil {
			if err := checkFile(filename); err != nil {
				return err
			}
		}

		// Determine the file path ensuring it's within targetDir.
		filePath := filepath.Join(targetDir,
			strings.ReplaceAll(filename, "/", string(os.PathSeparator)))

		// If it's a directory, create it.
		if f.FileInfo().IsDir() {
			log.Infofc(ctx, "create dir %s -> %s", f.Name, filePath)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create dir: %v", err)
			}
			continue
		}

		log.Infofc(ctx, "extract %s -> %s", f.Name, filePath)

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

			// log.Errorfc(ctx, "failed to copy file contents: %v", err)
			// continue
			return fmt.Errorf("failed to copy file contents: %v", err)
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
		if err != nil {
			return err
		}

		destPath := strings.TrimPrefix(path, srcDir)
		destPath = strings.ReplaceAll(destPath, `\`, "/")

		log.Infofc(ctx, "zipping %s...", path)
		if info.IsDir() {
			if destPath == "" {
				return nil
			}
			// ensure directory ends with /
			if !strings.HasSuffix(destPath, "/") {
				destPath += "/"
			}
			_, err := w.Create(destPath)
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		f, err := w.Create(destPath)
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
