package plateaucms

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/putil"
	"github.com/labstack/echo/v4"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/rerror"
	"github.com/samber/lo"
)

const (
	ProjectNameParam             = "pid"
	tokenProject                 = "system"
	metadataModel                = "workspaces"
	plateauProjectModel          = "plateau-projects"
	projectAliasField            = "project_alias"
	datacatalogProjectAliasField = "datacatalog_project_alias"
)

var HTTPMethodsAll = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
	http.MethodDelete,
}

var HTTPMethodsExceptGET = []string{
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
	http.MethodDelete,
}

type Config struct {
	CMSBaseURL      string
	CMSMainToken    string
	CMSTokenProject string
	// compat
	CMSMainProject string
	AdminToken     string
}

type CMS struct {
	cmsbase            string
	cmsMetadataProject string
	cmsMain            cms.Interface
	// comapt
	cmsMainProject string
	cmsToken       string
	adminToken     string
}

func New(c Config) (*CMS, error) {
	cmsMain, err := cms.New(c.CMSBaseURL, c.CMSMainToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cms: %w", err)
	}

	if c.CMSTokenProject == "" {
		c.CMSTokenProject = tokenProject
	}

	return &CMS{
		cmsbase:            c.CMSBaseURL,
		cmsMetadataProject: c.CMSTokenProject,
		cmsMain:            cmsMain,
		// compat
		cmsMainProject: c.CMSMainProject,
		cmsToken:       c.CMSMainToken,
		adminToken:     c.AdminToken,
	}, nil
}

func (h *CMS) Clone() *CMS {
	return &CMS{
		cmsbase:            h.cmsbase,
		cmsMetadataProject: h.cmsMetadataProject,
		cmsMain:            h.cmsMain,
		// compat
		cmsMainProject: h.cmsMainProject,
		cmsToken:       h.cmsToken,
		adminToken:     h.adminToken,
	}
}

func (h *CMS) AuthMiddleware(key string, authMethods []string, findDataCatalog bool, defaultProject string) echo.MiddlewareFunc {
	if key == "" {
		key = ProjectNameParam
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			prj := c.Param(key)
			if prj == "" {
				prj = defaultProject
			}

			md, all, err := h.Metadata(ctx, prj, findDataCatalog)
			if err != nil {
				if errors.Is(err, rerror.ErrNotFound) {
					ctx = context.WithValue(ctx, cmsMetadataContextKey{}, md)
					c.SetRequest(req.WithContext(ctx))
					return next(c)
				}
				return err
			}

			cmsh, err := cms.New(h.cmsbase, md.CMSAPIKey)
			if err != nil {
				return rerror.ErrInternalBy(fmt.Errorf("sidebar: failed to create cms for %s: %w", prj, err))
			}

			// auth
			header := req.Header.Get("Authorization")
			token := strings.TrimPrefix(header, "Bearer ")
			if md.SidebarAccessToken == "" || token != md.SidebarAccessToken {
				if len(authMethods) > 0 && slices.Contains(authMethods, req.Method) {
					return c.JSON(http.StatusUnauthorized, "unauthorized")
				}
			} else {
				md.Auth = true
			}

			// attach
			ctx = context.WithValue(ctx, cmsMetadataContextKey{}, md)
			ctx = context.WithValue(ctx, cmsAllMetadataContextKey{}, all)
			ctx = context.WithValue(ctx, cmsContextKey{}, cmsh)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

type cmsContextKey struct{}
type cmsMetadataContextKey struct{}
type cmsAllMetadataContextKey struct{}

func GetCMSFromContext(ctx context.Context) cms.Interface {
	cms, _ := ctx.Value(cmsContextKey{}).(cms.Interface)
	return cms
}

func GetCMSMetadataFromContext(ctx context.Context) Metadata {
	md, _ := ctx.Value(cmsMetadataContextKey{}).(Metadata)
	return md
}

func GetAllCMSMetadataFromContext(ctx context.Context) []Metadata {
	md, _ := ctx.Value(cmsAllMetadataContextKey{}).([]Metadata)
	return md
}

type Metadata struct {
	Name                     string `json:"name" cms:"name,text"`
	ProjectAlias             string `json:"project_alias" cms:"project_alias,text"`
	DataCatalogProjectAlias  string `json:"datacatalog_project_alias" cms:"datacatalog_project_alias,text"`
	DataCatalogSchemaVersion string `json:"datacatalog_schema_version" cms:"datacatalog_schema_version,select"`
	CMSAPIKey                string `json:"cms_apikey" cms:"cms_apikey,text"`
	SidebarAccessToken       string `json:"sidebar_access_token" cms:"sidebar_access_token,text"`
	SubPorjectAlias          string `json:"subproject_alias" cms:"subproject_alias,text"`
	// whether the request is authenticated with sidebar access token
	Auth       bool   `json:"-" cms:"-"`
	CMSBaseURL string `json:"-" cms:"-"`
}

func (h *CMS) Metadata(ctx context.Context, prj string, findDataCatalog bool) (Metadata, []Metadata, error) {
	// compat
	if h.cmsMainProject != "" && prj == h.cmsMainProject {
		return Metadata{
			ProjectAlias:       h.cmsMainProject,
			CMSAPIKey:          h.cmsToken,
			SidebarAccessToken: h.adminToken,
			CMSBaseURL:         h.cmsbase,
		}, nil, nil
	}

	all, err := h.AllMetadata(ctx, findDataCatalog)
	if err != nil {
		return Metadata{}, nil, err
	}

	md, ok := lo.Find(all, func(i Metadata) bool {
		return findDataCatalog && i.DataCatalogProjectAlias == prj || i.ProjectAlias == prj
	})
	if !ok {
		return Metadata{}, all, rerror.ErrNotFound
	}

	return md, all, nil
}

func (h *CMS) AllMetadata(ctx context.Context, findDataCatalog bool) ([]Metadata, error) {
	if h.cmsMetadataProject == "" {
		return nil, rerror.ErrNotFound
	}

	items, err := h.cmsMain.GetItemsByKeyInParallel(ctx, h.cmsMetadataProject, metadataModel, false, 100)
	if err != nil || items == nil {
		if errors.Is(err, cms.ErrNotFound) || items == nil {
			return nil, rerror.ErrNotFound
		}
		return nil, rerror.ErrInternalBy(fmt.Errorf("plateaucms: failed to get metadata: %w", err))
	}

	all := make([]Metadata, 0, len(items.Items))
	for _, item := range items.Items {
		m := Metadata{}
		item.Unmarshal(&m)
		if m.CMSAPIKey == "" {
			continue
		}
		if m.DataCatalogProjectAlias == "" {
			m.DataCatalogProjectAlias = m.ProjectAlias
		}
		m.CMSBaseURL = h.cmsbase
		all = append(all, m)
	}

	return all, nil
}

func (h *CMS) LastModified(c echo.Context, prj string, models ...string) (bool, error) {
	ctx := c.Request().Context()
	cmsh := GetCMSFromContext(ctx)

	mlastModified := time.Time{}
	for _, m := range models {
		model, err := cmsh.GetModelByKey(ctx, prj, m)
		if err != nil {
			if errors.Is(err, cms.ErrNotFound) {
				continue
			}
			return false, err
		}

		if model != nil && mlastModified.Before(model.LastModified) {
			mlastModified = model.LastModified
		}
	}

	return putil.LastModified(c, mlastModified)
}

func (m Metadata) CMS() (*cms.CMS, error) {
	return cms.New(m.CMSBaseURL, m.CMSAPIKey)
}
