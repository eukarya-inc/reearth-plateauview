package plateaucms

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/putil"
	"github.com/labstack/echo/v4"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/rerror"
	"github.com/samber/lo"
)

const (
	tokenProject      = "system"
	tokenModel        = "workspaces"
	tokenProjectField = "project_alias"
)

type Config struct {
	CMSBaseURL      string
	CMSMainToken    string
	CMSTokenProject string
	// compat
	CMSMainProject string
	AdminToken     string
}

type CMS struct {
	cmsbase         string
	cmsTokenProject string
	cmsMain         cms.Interface
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
		cmsbase:         c.CMSBaseURL,
		cmsTokenProject: c.CMSTokenProject,
		cmsMain:         cmsMain,
		// compat
		cmsMainProject: c.CMSMainProject,
		cmsToken:       c.CMSMainToken,
		adminToken:     c.AdminToken,
	}, nil
}

func (h *CMS) AuthMiddleware(skipAuth bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			prj := c.Param("pid")

			md, err := h.Metadata(ctx, prj)
			if err != nil {
				return err
			}

			cmsh, err := cms.New(h.cmsbase, md.CMSAPIKey)
			if err != nil {
				return rerror.ErrInternalBy(fmt.Errorf("sidebar: failed to create cms for %s: %w", prj, err))
			}

			// auth
			if !skipAuth && (req.Method == http.MethodPost || req.Method == http.MethodPatch || req.Method == http.MethodPut || req.Method == http.MethodDelete) {
				header := req.Header.Get("Authorization")
				token := strings.TrimPrefix(header, "Bearer ")
				if md.SidebarAccessToken == "" || token != md.SidebarAccessToken {
					return c.JSON(http.StatusUnauthorized, "unauthorized")
				}
			}

			// attach
			ctx = context.WithValue(ctx, cmsMetadataContextKey{}, md)
			ctx = context.WithValue(ctx, cmsContextKey{}, cmsh)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

type cmsContextKey struct{}
type cmsMetadataContextKey struct{}

func GetCMSFromContext(ctx context.Context) cms.Interface {
	return ctx.Value(cmsContextKey{}).(cms.Interface)
}

func GetCMSMetadataFromContext(ctx context.Context) Metadata {
	return ctx.Value(cmsMetadataContextKey{}).(Metadata)
}

type Metadata struct {
	ProjectAlias       string `json:"project_alias" cms:"project_alias,text"`
	CMSAPIKey          string `json:"cms_apikey" cms:"cms_apikey,text"`
	SidebarAccessToken string `json:"sidebar_access_token" cms:"sidebar_access_token,text"`
}

func (h *CMS) Metadata(ctx context.Context, prj string) (Metadata, error) {
	// compat
	if h.cmsMainProject != "" && prj == h.cmsMainProject {
		return Metadata{
			ProjectAlias:       h.cmsMainProject,
			CMSAPIKey:          h.cmsToken,
			SidebarAccessToken: h.adminToken,
		}, nil
	}

	if h.cmsTokenProject == "" {
		return Metadata{}, rerror.ErrNotFound
	}

	items, err := h.cmsMain.GetItemsByKeyInParallel(ctx, h.cmsTokenProject, tokenModel, false, 100)
	if err != nil || items == nil {
		if errors.Is(err, cms.ErrNotFound) || items == nil {
			return Metadata{}, rerror.ErrNotFound
		}
		return Metadata{}, rerror.ErrInternalBy(fmt.Errorf("sidebar: failed to get token: %w", err))
	}

	item, ok := lo.Find(items.Items, func(i cms.Item) bool {
		s := i.FieldByKey(tokenProjectField).ValueString()
		return s != nil && *s == prj
	})
	if !ok {
		return Metadata{}, rerror.ErrNotFound
	}

	m := Metadata{}
	item.Unmarshal(&m)
	if m.CMSAPIKey == "" {
		return Metadata{}, rerror.ErrNotFound
	}

	return m, nil
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
