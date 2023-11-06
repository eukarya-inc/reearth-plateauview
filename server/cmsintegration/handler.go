package cmsintegration

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationv2"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationv3"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/dataconv"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/geospatialjp"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

type Config = cmsintegrationcommon.Config

func Handler(conf Config, g *echo.Group) error {
	// v3
	v3, err := cmsintegrationv3.NotifyHandler(conf)
	if err != nil {
		return err
	}

	g.POST("/notify_fme/v3", v3)

	// v2 (compat)
	return compatHandler(conf, g)
}

func compatHandler(conf Config, g *echo.Group) error {
	v2, err := cmsintegrationv2.NotifyHandler(conf)
	if err != nil {
		return err
	}

	geo, err := geospatialjp.Handler(geospatialjpConfig(conf))
	if err != nil {
		return err
	}

	dataconv, err := dataconv.Handler(dataConvConfig(conf))
	if err != nil {
		return err
	}

	g.POST("/notify_fme", v2)
	g.POST("/publish_to_geospatialjp", geo)
	g.POST("/dataconv", echo.WrapHandler(dataconv))
	return nil
}

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	h1, err := cmsintegrationv3.WebhookHandler(conf)
	if err != nil {
		return nil, err
	}

	h2, err := cmsintegrationv2.WebhookHandler(conf)
	if err != nil {
		return nil, err
	}

	h3, err := geospatialjp.WebhookHandler(geospatialjpConfig(conf))
	if err != nil {
		return nil, err
	}

	h4, err := dataconv.WebhookHandler(dataConvConfig(conf))
	if err != nil {
		return nil, err
	}

	return cmswebhook.MergeHandlers([]cmswebhook.Handler{
		h1, h2, h3, h4,
	}), nil
}

func geospatialjpConfig(conf Config) geospatialjp.Config {
	return geospatialjp.Config{
		CMSBase:             conf.CMSBaseURL,
		CMSToken:            conf.CMSToken,
		CMSIntegration:      conf.CMSIntegration,
		CkanBase:            conf.CkanBaseURL,
		CkanOrg:             conf.CkanOrg,
		CkanToken:           conf.CkanToken,
		CkanPrivate:         conf.CkanPrivate,
		DisablePublication:  conf.DisableGeospatialjpPublication,
		DisableCatalogCheck: conf.DisableGeospatialjpCatalogCheck,
		PublicationToken:    conf.APIToken,
		// EnablePulicationOnWebhook: true,
	}
}

func dataConvConfig(conf Config) dataconv.Config {
	return dataconv.Config{
		Disable:  conf.DisableDataConv,
		CMSBase:  conf.CMSBaseURL,
		CMSToken: conf.CMSToken,
		APIToken: conf.APIToken,
		// CMSModel: conf.CMSModel,
	}
}
