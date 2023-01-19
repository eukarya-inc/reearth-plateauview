package sdkapi

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/rerror"
)

func Handler(conf Config, g *echo.Group) {
	conf.Normalize()
	cl := NewClient(nil, conf.CMSBaseURL, conf.Project, conf.Model)

	g.GET("/datasets", func(c echo.Context) error {
		data, err := Datasets(c.Request().Context(), cl)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})

	g.GET("/datasets/:id/files", func(c echo.Context) error {
		data, err := Files(c.Request().Context(), cl, c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})
}

func Datasets(ctx context.Context, c *Client) (*DatasetResponse, error) {
	items, err := c.GetItems(ctx)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return items.DatasetResponse(), nil
}

func Files(ctx context.Context, c *Client, id string) (any, error) {
	item, err := c.GetItem(ctx, id)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}
	if item.CityGML == nil || item.MaxLOD == nil {
		return nil, rerror.ErrNotFound
	}

	maxlod, err := c.GetMaxLOD(item.MaxLOD.URL)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return maxlod.Map().Files(item.CityGML.URL), nil
}
