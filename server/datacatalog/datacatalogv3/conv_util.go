package datacatalogv3

import (
	"path"
	"strings"

	"github.com/samber/lo"
)

func namesWithoutExtFromUrls(urls []string) []string {
	res := make([]string, 0, len(urls))
	for _, url := range urls {
		res = append(res, nameWithoutExt(nameFromUrl(url)))
	}
	return res
}

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

func dropNil[T any](s []*T) []*T {
	res := make([]*T, 0, len(s))
	for _, i := range s {
		if i != nil {
			res = append(res, i)
		}
	}
	return res
}
