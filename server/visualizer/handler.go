package visualizer

import (
	"errors"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// 仮置
var (
	fieldId         = "01gkjdq9h2t4c4x300fewj98zq"
	dataModelId     = "01gkjdpwkh478tysj1xc4wj3cc"
	templateModelId = ""
)

// GET | /viz/:id
func fetchRoot(CMS cms.Interface) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// projectIDが必要な場合使う
		_ = c.Param("pid")
		data, err := CMS.GetItems(ctx, dataModelId) // modelID: "plateau-view-data"
		if err != nil {
			return err
		}
		templates, err := CMS.GetItems(ctx, templateModelId) // modelID: "templates"
		if err != nil {
			return err
		}
		root := ToRoot(templates, data)
		return c.JSON(200, &root)
	}
}

// GET | /viz/:pid/data
func getDataHandler(CMS cms.Interface) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//projectIDが必要な場合使う
		_ = c.Param("pid")
		var data any
		if err := c.Bind(&data); err != nil {
			//TODO: エラーハンドリングをきれいにする
			return fmt.Errorf("failed to bind a data: %w", err)
		}

		fields := []cms.Field{{
			ID:    fieldId,
			Value: data,
		}}
		item, err := CMS.CreateItem(ctx, dataModelId, fields) //modelID: "plateau-view-data"
		if err != nil {
			return err
		}

		field, found := lo.Find(item.Fields, func(i cms.Field) bool {
			return fieldId == i.ID
		})
		if !found {
			//TODO: エラーハンドリングをきれいにする
			return errors.New("err")
		}

		res := Component{
			ID:        item.ID,
			Component: field.Value,
		}
		return c.JSON(200, &res)
	}
}

// POST | /viz/:pid/data
func createDataHandler(CMS cms.Interface) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//projectIDが必要な場合使う
		_ = c.Param("pid")
		var data any
		if err := c.Bind(&data); err != nil {
			//TODO: エラーハンドリングをきれいにする
			return fmt.Errorf("failed to bind a data: %w", err)
		}

		fields := []cms.Field{{
			ID:    fieldId,
			Value: data,
		}}
		item, err := CMS.CreateItem(ctx, dataModelId, fields) //modelID: "plateau-view-data"
		if err != nil {
			return err
		}

		field, found := lo.Find(item.Fields, func(i cms.Field) bool {
			return fieldId == i.ID
		})
		if !found {
			//TODO: エラーハンドリングをきれいにする
			return errors.New("err")
		}

		res := Component{
			ID:        item.ID,
			Component: field.Value,
		}
		return c.JSON(200, &res)
	}
}

// PATCH | /viz/:pid/data/:did
func updateDataHandler(CMS cms.Interface) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//projectIDが必要な場合使う
		_ = c.Param("pid")
		var data any
		if err := c.Bind(&data); err != nil {
			//TODO: エラーハンドリングをきれいにする
			return fmt.Errorf("failed to bind a data: %w", err)
		}

		fields := []cms.Field{{
			ID:    fieldId,
			Value: data,
		}}
		item, err := CMS.UpdateItem(ctx, dataModelId, fields) //modelID: "plateau-view-data"
		if err != nil {
			return err
		}

		field, found := lo.Find(item.Fields, func(i cms.Field) bool {
			return fieldId == i.ID
		})
		if !found {
			//TODO: エラーハンドリングをきれいにする
			return fmt.Errorf("not found a fields data: %w", err)
		}

		res := Component{
			ID:        item.ID,
			Component: field.Value,
		}
		return c.JSON(200, &res)
	}
}

// POST | /viz/:pid/templates
func createTemplateHandler(CMS cms.Interface) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//0. プロジェクトのIDを取得する
		_ = c.Param("pid")
		//1. FEから来たJSONを取得する
		var data any
		if err := c.Bind(&data); err != nil {
			//TODO: エラーハンドリングをきれいにする
			return fmt.Errorf("failed to bind a data: %w", err)
		}
		//2. JSONをCMSに登録する
		//2-1. CMSにわたすデータを作成する
		fields := []cms.Field{}
		//2-2. CMSにわたす
		//TODO: モデルのIDを後で設定から読み込むように変更する
		item, err := CMS.CreateItem(ctx, "template", fields)
		if err != nil {
			return err
		}
		//3. 結果を返す
		res := Component{
			ID:        item.ID,
			Component: item.Fields[0].Value,
		}
		return c.JSON(200, &res)
	}
}
