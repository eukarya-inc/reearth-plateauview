package searchindex

import (
	"net/url"
	"path"
)

type Config struct {
	CMSBase           string
	CMSToken          string
	CMSStorageProject string
	// optioanl
	CMSStorageModel string
	// internal
	skipIndexer bool
}

func getAssetBase(u *url.URL) string {
	u2 := *u
	b := path.Join(path.Dir(u.Path), pathFileName(u.Path))
	u2.Path = b
	return u2.String()
}
