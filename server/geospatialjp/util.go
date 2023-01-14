package geospatialjp

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/geospatialjp/ckan"
)

func catalogFinalFileName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	b := path.Base(u.Path)

	be, af, found := strings.Cut(b, ".")
	if !found {
		return "", fmt.Errorf("invalid file name: %s", b)
	}

	return fmt.Sprintf("%s_final.%s", be, af), nil
}

func resources(pkg *ckan.Package, c Catalog, citygmlURL, allURL, catalogURL string, private bool) []ckan.Resource {
	// TODO
	return nil
}

func packageKey(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	b := path.Base(u.Path)

	// TODO

	return b, nil
}

func packageFromCatalog(c Catalog, pkgID string) ckan.Package {
	return ckan.Package{
		ID: pkgID,
		// TODO
	}
}
