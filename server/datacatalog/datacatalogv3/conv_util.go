package datacatalogv3

import (
	"path"
	"strings"

	"github.com/samber/lo"
)

func nameFromUrl(url string) string {
	if url == "" {
		return ""
	}

	if i := strings.LastIndexByte(url, '/'); i >= 0 {
		url = url[i+1:]
	}

	return url
}

func nameWithoutExt(name string) string {
	ext := path.Ext(name)
	if ext == "" {
		return name
	}

	return name[:len(name)-len(ext)]
}

func toPtrIfPresent[T comparable](v T) *T {
	if lo.IsEmpty(v) {
		return nil
	}
	return &v
}

func firstNonEmptyValue[T comparable](v ...T) (_ T) {
	for _, i := range v {
		if !lo.IsEmpty(i) {
			return i
		}
	}
	return
}
