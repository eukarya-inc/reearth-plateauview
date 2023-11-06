package cmsintegrationv2

import (
	"errors"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

type Config = cmsintegrationcommon.Config

type Services struct {
	FME fmeInterface
	CMS cms.Interface
}

func NewServices(c Config) (s *Services, _ error) {
	s = &Services{}

	if !c.FMEMock {
		fmeURL := c.FMEURLV3
		if fmeURL == "" {
			return nil, errors.New("FME URL is not set")
		}
		fme := newFME(fmeURL, c.FMEResultURL, c.FMESkipQualityCheck)
		s.FME = fme
	}

	cms, err := cms.New(c.CMSBaseURL, c.CMSToken)
	if err != nil {
		return nil, fmt.Errorf("failed to init cms: %w", err)
	}
	s.CMS = cms

	return
}
