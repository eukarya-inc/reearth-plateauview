package main

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog"
	"github.com/eukarya-inc/reearth-plateauview/server/opinion"
	"github.com/eukarya-inc/reearth-plateauview/server/sdkapi"
	"github.com/eukarya-inc/reearth-plateauview/server/searchindex"
	"github.com/eukarya-inc/reearth-plateauview/server/sidebar"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/util"
)

type Service struct {
	Name           string
	Echo           func(g *echo.Group) error
	Webhook        cmswebhook.Handler
	DisableNoCache bool
}

var services = [](func(*Config) (*Service, error)){
	CMSIntegration,
	SDKAPI,
	SearchIndex,
	Opinion,
	Sidebar,
	DataCatalog,
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
	if c.CMSBaseURL == "" || c.CMSToken == "" || c.FMEBaseURL == "" || c.Host == "" || c.FMEToken == "" {
		return nil, nil
	}

	w, err := cmsintegration.WebhookHandler(c)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name: "cmsintegration",
		Echo: func(g *echo.Group) error {
			return cmsintegration.Handler(c, g)
		},
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

func SDKAPI(conf *Config) (*Service, error) {
	c := conf.SDKAPI()
	if c.CMSBaseURL == "" || c.Project == "" {
		return nil, nil
	}

	return &Service{
		Name:           "sdkapi",
		DisableNoCache: true,
		Echo: func(g *echo.Group) error {
			return sdkapi.Handler(c, g.Group("/sdk"))
		},
	}, nil
}

func Opinion(conf *Config) (*Service, error) {
	c := conf.Opinion()
	if c.SendGridAPIKey == "" || c.From == "" || c.To == "" {
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

func Sidebar(conf *Config) (*Service, error) {
	c := conf.Sidebar()
	if c.AdminToken == "" || c.CMSMainToken == "" || c.CMSBaseURL == "" {
		return nil, nil
	}

	return &Service{
		Name:           "sidebar",
		DisableNoCache: true,
		Echo: func(g *echo.Group) error {
			return util.Try(
				func() error { return sidebar.Echo(g.Group("/sidebar"), c) },
				func() error { return sidebar.ShareEcho(g.Group("/share"), c) },
			)
		},
	}, nil
}

func DataCatalog(conf *Config) (*Service, error) {
	c := conf.DataCatalog()
	if c.CMSBase == "" {
		return nil, nil
	}
	if c.PlaygroundEndpoint == "" {
		c.PlaygroundEndpoint = "/datacatalog"
	}

	return &Service{
		Name: "datacatalog",
		Echo: func(g *echo.Group) error {
			return datacatalog.Echo(c, g.Group("/datacatalog"))
		},
		DisableNoCache: true,
	}, nil
}
