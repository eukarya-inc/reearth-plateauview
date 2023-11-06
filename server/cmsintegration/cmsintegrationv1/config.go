package cmsintegrationv1

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

type Config = cmsintegrationcommon.Config

type Services struct {
	FME fmeInterface
	CMS cms.Interface
}

func NewServices(c Config) (s Services, _ error) {
	if !c.FMEMock {
		fme, err := NewFME(c.FMEBaseURL, c.FMEToken, c.FMEResultURL)
		if err != nil {
			return Services{}, fmt.Errorf("failed to init fme: %w", err)
		}
		s.FME = fme
	}

	cms, err := cms.New(c.CMSBaseURL, c.CMSToken)
	if err != nil {
		return Services{}, fmt.Errorf("failed to init cms: %w", err)
	}
	s.CMS = cms

	return
}
