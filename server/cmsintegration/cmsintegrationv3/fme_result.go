package cmsintegrationv3

import (
	"fmt"
	"sort"
	"strings"
)

const fmeIDPrefix = "v3"

type fmeID struct {
	ItemID      string
	ProjectID   string
	FeatureType string
}

func parseFMEID(id, secret string) (fmeID, error) {
	payload, err := unsignFMEID(id, secret)
	if err != nil {
		return fmeID{}, err
	}

	s := strings.SplitN(payload, ";", 4)
	if len(s) != 4 || s[0] != fmeIDPrefix {
		return fmeID{}, ErrInvalidFMEID
	}

	return fmeID{
		ItemID:      s[1],
		ProjectID:   s[2],
		FeatureType: s[3],
	}, nil
}

func (i fmeID) String(secret string) string {
	payload := fmt.Sprintf("%s;%s;%s;%s", fmeIDPrefix, i.ItemID, i.ProjectID, i.FeatureType)
	return signFMEID(payload, secret)
}

type fmeResult struct {
	Type    string         `json:"type"`
	Status  string         `json:"status"`
	ID      string         `json:"id"`
	LogURL  string         `json:"logUrl"`
	Results map[string]any `json:"results"`
}

func (f fmeResult) ParseID(secret string) fmeID {
	id, err := parseFMEID(f.ID, secret)
	if err != nil {
		return fmeID{}
	}
	return id
}

type fmeResultURLs struct {
	FeatureType string
	Data        []string
	Dic         string
	MaxLOD      string
}

func (f fmeResult) GetResultURLs(featureType string) (res fmeResultURLs) {
	res.FeatureType = featureType

	for k, v := range f.Results {
		if k == featureType || strings.HasPrefix(k, featureType+"/") || strings.HasPrefix(k, featureType+"_") {
			if v2, ok := v.(string); ok {
				res.Data = append(res.Data, v2)
			} else if v2, ok := v.([]any); ok {
				for _, v3 := range v2 {
					if v4, ok := v3.(string); ok {
						res.Data = append(res.Data, v4)
					}
				}
			}
		}
	}

	sort.Strings(res.Data)

	if v, ok := f.Results["_dic"].(string); ok {
		res.Dic = v
	}

	if v, ok := f.Results["_maxlod"].(string); ok {
		res.MaxLOD = v
	}

	return
}
