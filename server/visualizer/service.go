package visualizer

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/fme"
)

type Config struct {
	CMSModelID       string
	CMSBaseURL       string
	CMSToken         string
	VToken           string
	DataModelKey     string
	TemplateModelKey string
}

type Services struct {
	FME fme.Interface
	CMS cms.Interface
}

func NewServices(c Config) (s Services, _ error) {
	cms, err := cms.New(c.CMSBaseURL, c.CMSToken)
	if err != nil {
		return Services{}, fmt.Errorf("failed to init cms: %w", err)
	}
	s.CMS = cms

	return
}
