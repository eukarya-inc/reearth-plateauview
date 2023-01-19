package visualizer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Handler struct {
	DataModelKey                 string
	TemplateModelKey             string
	DataModelDataFieldID         string
	DataModelIDFieldID           string
	TemplateModelTemplateFieldID string
	TemplateModelIDFieldID       string
	CMS                          cms.Interface
}

func NewHandler(CMS cms.Interface, dKey, tKey string) (*Handler, error) {
	if dKey == "" || tKey == "" {
		return nil, fmt.Errorf("missing model keys, dataKey=%s, templateKey=%s", dKey, tKey)
	}

	ctx := context.Background()
	data, err := CMS.GetItems(ctx, dKey)
	if err != nil {
		return nil, err
	}

	templates, err := CMS.GetItems(ctx, tKey)
	if err != nil {
		return nil, err
	}

	if len(data.Items) == 0 || len(templates.Items) == 0 {
		return nil, fmt.Errorf("failed to fetch meta data")
	}

	h := &Handler{
		DataModelKey:                 dKey,
		TemplateModelKey:             tKey,
		CMS:                          CMS,
		DataModelIDFieldID:           data.Items[0].Fields[0].ID,
		DataModelDataFieldID:         data.Items[0].Fields[1].ID,
		TemplateModelIDFieldID:       templates.Items[0].Fields[0].ID,
		TemplateModelTemplateFieldID: templates.Items[0].Fields[1].ID,
	}

	return h, nil
}

// GET | /viz/:id
func (h *Handler) fetchRoot() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		data, err := h.CMS.GetItems(ctx, h.DataModelKey)
		if err != nil {
			return err
		}

		data2 := lo.Map(data.Items, func(d cms.Item, _ int) any {
			return d.Field(h.DataModelDataFieldID).Value
		})

		templates, err := h.CMS.GetItems(ctx, h.TemplateModelKey)
		if err != nil {
			return err
		}

		templates2 := lo.Map(templates.Items, func(t cms.Item, _ int) any {
			return t.Field(h.TemplateModelTemplateFieldID).Value
		})

		root := struct {
			Templates  []any `json:"templates"`
			Components []any `json:"data"`
		}{
			Templates:  templates2,
			Components: data2,
		}

		return c.JSON(http.StatusOK, root)
	}
}

// GET | /viz/:pid/data
func (h *Handler) getAllDataHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// res, err := http.Get(lo.Must(url.JoinPath("https://api.cms.test.reearth.dev", "api/models/key1/items")))
		// fmt.Printf(res.Status)
		data, err := h.CMS.GetItems(ctx, h.DataModelKey)
		if err != nil {
			return err
		}

		data2 := lo.Map(data.Items, func(d cms.Item, _ int) any {
			return d.Field(h.DataModelDataFieldID)
		})
		return c.JSON(http.StatusOK, data2)
	}
}

// GET | /viz/:pid/data/:iid
func (h *Handler) getDataHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		itemID := c.Param("iid")
		if itemID == "" {
			return c.JSON(http.StatusNotFound, nil)
		}

		data, err := h.CMS.GetItem(ctx, itemID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data.Field(h.DataModelIDFieldID).Value)
	}
}

// POST | /viz/:pid/data
func (h *Handler) createDataHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("occur an unexpected EOF error: %w", err)
		}

		fields := []cms.Field{{
			ID:    h.DataModelDataFieldID,
			Value: string(b),
		}}
		item, err := h.CMS.CreateItem(ctx, h.DataModelKey, fields)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, item.Field(h.DataModelDataFieldID).Value)
	}
}

// PATCH | /viz/:pid/data/:did
func (h *Handler) updateDataHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		itemID := c.Param("iid")
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("occur an unexpected EOF error: %w", err)
		}

		fields := []cms.Field{{
			ID:    h.DataModelDataFieldID,
			Value: string(b),
		}}

		item, err := h.CMS.UpdateItem(ctx, itemID, fields)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, item.Field(h.DataModelDataFieldID).Value)
	}
}

// DELETE | /viz/:pid/data/:did
func (h *Handler) deleteDataHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		itemID := c.Param("iid")

		err := h.CMS.DeleteItem(ctx, itemID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// GET | /viz/:id/templates
func (h *Handler) fetchTemplate() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		templates, err := h.CMS.GetItems(ctx, h.TemplateModelKey)
		if err != nil {
			return err
		}

		templates2 := lo.Map(templates.Items, func(t cms.Item, _ int) any {
			return t.Field(h.TemplateModelTemplateFieldID).Value
		})
		return c.JSON(http.StatusOK, templates2)
	}
}

// POST | /viz/:pid/templates
func (h *Handler) createTemplateHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("occur an unexpected EOF error: %w", err)
		}

		fields := []cms.Field{{
			ID:    h.TemplateModelTemplateFieldID,
			Value: string(b),
		}}
		item, err := h.CMS.CreateItem(ctx, h.TemplateModelKey, fields)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, item.Field(h.TemplateModelTemplateFieldID))
	}
}

// PATCH | /viz/:id/templates/:itemId
func (h *Handler) updateTemplateHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		itemID := c.Param("iid")
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("occur an unexpected EOF error: %w", err)
		}

		fields := []cms.Field{{
			ID:    h.TemplateModelTemplateFieldID,
			Value: string(b),
		}}

		item, err := h.CMS.UpdateItem(ctx, itemID, fields)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, item.Field(h.TemplateModelTemplateFieldID).Value)
	}
}

// DELETE | /viz/:id/templates/:itemId
func (h *Handler) deleteTemplateHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		itemID := c.Param("iid")

		err := h.CMS.DeleteItem(ctx, itemID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, nil)
	}
}
