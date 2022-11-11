package cmsintegration

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/fme"
)

type Config struct {
	FMEBaseURL       string
	FMEToken         string
	FMEResultURL     string
	CMSBaseURL       string
	CMSToken         string
	CMSWebhookSecret string
	Secret           string
}

type Services struct {
	FME fme.Interface
	CMS cms.Interface
}

func NewServices(c Config) (Services, error) {
	fme, err := fme.New(c.FMEBaseURL, c.FMEToken, c.FMEResultURL+"/notify")
	if err != nil {
		return Services{}, fmt.Errorf("failed to init fme: %w", err)
	}

	cms, err := cms.New(c.CMSBaseURL, c.CMSToken)
	if err != nil {
		return Services{}, fmt.Errorf("failed to init cms: %w", err)
	}

	return Services{FME: fme, CMS: cms}, nil
}
