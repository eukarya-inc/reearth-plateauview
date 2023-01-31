package visualizer

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
)

type Config struct {
	CMSProject       string
	CMSModelID       string
	CMSBaseURL       string
	CMSToken         string
	AdminToken       string
	DataModelKey     string
	TemplateModelKey string
}

type Services struct {
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
