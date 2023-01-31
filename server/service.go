package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cms/cmswebhook"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration"
	"github.com/eukarya-inc/reearth-plateauview/server/geospatialjp"
	"github.com/eukarya-inc/reearth-plateauview/server/opinion"
	"github.com/eukarya-inc/reearth-plateauview/server/sdk"
	"github.com/eukarya-inc/reearth-plateauview/server/sdkapi"
	"github.com/eukarya-inc/reearth-plateauview/server/searchindex"
	"github.com/eukarya-inc/reearth-plateauview/server/share"
	"github.com/eukarya-inc/reearth-plateauview/server/visualizer"
	"github.com/labstack/echo/v4"
)

type Service struct {
	Name    string
	Echo    func(g *echo.Group) error
	Webhook cmswebhook.Handler
}

var services = [](func(*Config) (*Service, error)){
	CMSIntegration,
	Geospatialjp,
	SDK,
	SDKAPI,
	SearchIndex,
	Share,
	Opinion,
	SidebarAPI,
}

func Services(conf *Config) (srv []*Service, _ error) {
	for _, i := range services {
		s, err := i(conf)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", funcName(i), err)
		}
		if s == nil {
			continue
		}
		srv = append(srv, s)
	}
	return
}

func CMSIntegration(conf *Config) (*Service, error) {
	c := conf.CMSIntegration()
	if c.CMSBaseURL == "" || c.CMSToken == "" || c.FMEBaseURL == "" || c.FMEResultURL == "" || c.FMEToken == "" {
		return nil, nil
	}

	e, err := cmsintegration.NotifyHandler(c)
	if err != nil {
		return nil, err
	}

	w, err := cmsintegration.WebhookHandler(c)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name: "cmsintegration",
		Echo: func(g *echo.Group) error {
			g.POST("/notify_fme", e)
			return nil
		},
		Webhook: w,
	}, nil
}

func Geospatialjp(conf *Config) (*Service, error) {
	c := conf.Geospatialjp()
	if c.CMSBase == "" || c.CMSToken == "" || c.CkanBase == "" || c.CkanToken == "" || c.CkanOrg == "" {
		return nil, nil
	}

	w, err := geospatialjp.WebhookHandler(c)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:    "geospatialjp",
		Webhook: w,
	}, nil
}

func SearchIndex(conf *Config) (*Service, error) {
	c := conf.SearchIndex()
	if c.CMSBase == "" || c.CMSToken == "" || c.CMSStorageProject == "" {
		return nil, nil
	}

	w, err := searchindex.WebhookHandler(c)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:    "searchindex",
		Webhook: w,
	}, nil
}

func SDK(conf *Config) (*Service, error) {
	c := conf.SDK()
	if c.CMSBase == "" || c.CMSToken == "" || c.FMEBaseURL == "" || c.FMEResultURL == "" || c.FMEToken == "" {
		return nil, nil
	}

	e, err := sdk.NotifyHandler(c)
	if err != nil {
		return nil, err
	}

	w, err := sdk.WebhookHandler(c)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name: "sdk",
		Echo: func(g *echo.Group) error {
			g.POST("/notify_sdk", e)
			return nil
		},
		Webhook: w,
	}, nil
}

func SDKAPI(conf *Config) (*Service, error) {
	c := conf.SDKAPI()
	if c.CMSBaseURL == "" || c.Project == "" {
		return nil, nil
	}

	return &Service{
		Name: "sdkapi",
		Echo: func(g *echo.Group) error {
			sdkapi.Handler(c, g.Group("/sdk"))
			return nil
		},
	}, nil
}

func Share(conf *Config) (*Service, error) {
	c := conf.Share()
	if c.CMSBase == "" || c.CMSToken == "" || c.CMSProject == "" {
		return nil, nil
	}

	return &Service{
		Name: "share",
		Echo: func(g *echo.Group) error {
			return share.Echo(g.Group("/share"), c)
		},
	}, nil
}

func Opinion(conf *Config) (*Service, error) {
	c := conf.Opinion()
	if c.SendGridAPIKey == "" {
		return nil, nil
	}

	return &Service{
		Name: "opinion",
		Echo: func(g *echo.Group) error {
			opinion.Echo(g.Group("/opinion"), c)
			return nil
		},
	}, nil
}

func SidebarAPI(conf *Config) (*Service, error) {
	c := conf.Visualizer()
	if c.AdminToken == "" || c.CMSToken == "" {
		return nil, nil
	}

	return &Service{
		Name: "sidebar",
		Echo: func(g *echo.Group) error {
			return visualizer.Echo(g.Group("/visualizer"), c)
		},
	}, nil
}

func funcName(i interface{}) string {
	return strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), "main.")
}
