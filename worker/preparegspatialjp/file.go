package preparegspatialjp

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

func downloadAndUnzip(ctx context.Context, url, dir, tmpdir string, options *UnzipOptions) (string, error) {
	r, le, p, err := downloadFileAsReaderAtCloser(ctx, url, tmpdir)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = r.Close()
	}()

	return p, Unzip(ctx, r, le, dir, options)
}

func downloadFileAsReaderAtCloser(ctx context.Context, url, dir string) (ReaderAtCloser, int64, string, error) {
	p, err := downloadFileTo(ctx, url, dir)
	if err != nil {
		return nil, 0, "", err
	}

	stat, err := os.Stat(p)
	if err != nil {
		return nil, 0, "", err
	}

	b, err := os.Open(p)
	if err != nil {
		return nil, 0, "", err
	}

	return b, stat.Size(), p, nil
}

func downloadAndConsumeZip(ctx context.Context, url, dir string, fn func(*zip.Reader, os.FileInfo) error) error {
	return downloadAndConsumeFile(ctx, url, dir, func(f *os.File, fi os.FileInfo) error {
		zr, err := zip.NewReader(f, fi.Size())
		if err != nil {
			return err
		}
		return fn(zr, fi)
	})
}

func downloadAndConsumeFile(ctx context.Context, url, dir string, fn func(f *os.File, fi os.FileInfo) error) error {
	p, err := downloadFileTo(ctx, url, dir)
	if err != nil {
		return err
	}

	return consumeFile(p, fn)
}

func downloadFileTo(ctx context.Context, url, dir string) (string, error) {
	r, err := downloadFile(ctx, url)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = r.Close()
	}()

	name := fileNameFromURL(url)
	dest := filepath.Join(dir, name)
	_ = os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = io.Copy(f, r)
	if err != nil {
		return "", err
	}
	return dest, nil
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
}

func consumeFile(p string, fn func(f *os.File, fi os.FileInfo) error) (err error) {
	s, err2 := os.Stat(p)
	if err2 != nil {
		err = err2
		return
	}

	f, err2 := os.Open(p)
	if err2 != nil {
		err = err2
		return
	}

	defer func() {
		_ = f.Close()
		if err == nil {
			_ = os.Remove(p)
		}
	}()

	err = fn(f, s)
	return
}

type UnzipOptions struct {
	Prefix string
	Rename func(string) (string, error)
}

var SkipUnzip = fmt.Errorf("skip unzip")

func Unzip(ctx context.Context, zipFile io.ReaderAt, le int64, targetDir string, options *UnzipOptions) error {
	opts := lo.FromPtr(options)

	// Open the zip archive for reading.
	r, err := zip.NewReader(zipFile, le)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}

	// Iterate through the files in the archive.
	for _, f := range r.File {
		filename := f.Name
		filename = strings.ReplaceAll(filename, `\`, "/")
		filename = strings.TrimSuffix(filename, "/")

		if opts.Prefix != "" {
			skipped := false
			if strings.HasPrefix(filename, opts.Prefix) {
				filename = strings.TrimPrefix(filename, opts.Prefix)
			} else {
				skipped = true
			}

			// if prefix is "xxx/", xxx/hoge is skipped, and xxx should be also skipped
			if !skipped && len(opts.Prefix) > 1 && strings.HasSuffix(opts.Prefix, "/") && filename == opts.Prefix[:len(opts.Prefix)-1] {
				skipped = true
			}

			if skipped {
				log.Debugf("skip %s (no prefix: %s)", f.Name, opts.Prefix)
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

		if opts.Rename != nil {
			if fn, err := opts.Rename(filename); err != nil && err != SkipUnzip {
				return err
			} else if err == SkipUnzip {
				log.Debugf("skip %s", f.Name)
				continue
			} else if fn != "" {
				filename = fn
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

func ZipDir(ctx context.Context, src, dest string, destPrefix bool) error {
	log.Infofc(ctx, "start zipping %s to %s...", src, dest)

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		destPath := strings.TrimPrefix(p, src)
		destPath = strings.ReplaceAll(destPath, `\`, "/")
		if destPrefix {
			base := filepath.Base(src)
			destPath = path.Join(base, destPath)
		}

		log.Infofc(ctx, "zipping %s -> %s ...", p, destPath)
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

		file, err := os.Open(p)
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

	if err := filepath.Walk(src, walker); err != nil {
		return fmt.Errorf("failed to walk dir: %v", err)
	}

	return nil
}

func normalizeZipFilePath(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	if strings.HasPrefix(p, "__MACOSX/") ||
		strings.HasSuffix(p, "/.DS_Store") ||
		strings.HasSuffix(p, "/Thumb.db") ||
		p == ".DS_Store" || p == "Thumbs.db" {
		return ""
	}
	return p
}
