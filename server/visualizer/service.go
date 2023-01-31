package visualizer

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration"
	"github.com/eukarya-inc/reearth-plateauview/server/opinion"
	"github.com/eukarya-inc/reearth-plateauview/server/sdk"
	"github.com/eukarya-inc/reearth-plateauview/server/share"
	"github.com/labstack/echo"
	"github.com/samber/lo"
)

type Config struct {
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

	e := echo.New()
	e.POST("/notify_fme", lo.Must(cmsintegration.NotifyHandler(conf.CMSIntegration())))
	e.POST("/notify_sdk", lo.Must(sdk.NotifyHandler(conf.SDK())))
	lo.Must0(visualizer.Echo(e.Group(""), conf.Visualizer()))
	lo.Must0(share.Echo(e.Group("/share"), conf.Share()))
	opinion.Echo(e.Group("/opinion"), conf.Opinion())

	return
}
