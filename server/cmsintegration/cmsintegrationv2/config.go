package cmsintegrationv2

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	"github.com/eukarya-inc/reearth-plateauview/server/fme"
	cms "github.com/reearth/reearth-cms-api/go"
)

type Config = cmsintegrationcommon.Config

type Services struct {
	FME fme.Interface
	CMS cms.Interface
}

func NewServices(c Config) (s Services, _ error) {
	if !c.FMEMock {
		fme, err := fme.New(c.FMEBaseURL, c.FMEToken, c.FMEResultURL)
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
