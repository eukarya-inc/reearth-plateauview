package sdk

import (
	"fmt"
	"net/url"
	"strings"
)

type fmeID struct {
	ItemID    string
	AssetID   string
	ProjectID string
}

func parseFMEID(id, secret string) (fmeID, error) {
	payload, err := unsignFMEID(id, secret)
	if err != nil {
		return fmeID{}, err
	}

	s := strings.SplitN(payload, ";", 4)
	if len(s) != 3 {
		return fmeID{}, ErrInvalidFMEID
	}

	return fmeID{
		ItemID:    s[0],
		AssetID:   s[1],
		ProjectID: s[2],
	}, nil
}

func (i fmeID) String(secret string) string {
	payload := fmt.Sprintf("%s;%s;%s", i.ItemID, i.AssetID, i.ProjectID)
	return signFMEID(payload, secret)
}

type maxLODRequest struct {
	ID     string
	Target string
}

func (r maxLODRequest) Query() url.Values {
	q := url.Values{}
	q.Set("id", r.ID)
	q.Set("url", r.Target)
	return q
}

func (r maxLODRequest) Name() string {
	return "plateau2022-cms/maxlod-extract"
}
