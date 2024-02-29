package preparegspatialjp

import (
	"context"
	"io"
	"os"
	"path/filepath"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type CMSWrapper struct {
	CMS         cms.Interface
	ProjectID   string
	DataItemID  string
	CityItemID  string
	SkipCityGML bool
	SkipPlateau bool
	SkipMaxLOD  bool
	SkipIndex   bool
	SkipRelated bool
}

func (c *CMSWrapper) NotifyRunning(ctx context.Context, citygml, plateau, maxlod bool) {
	if c == nil {
		return
	}

	c.Comment(ctx, "公開準備処理を開始しました。")

	item := &GspatialjpDataItem{}
	if citygml {
		item.MergeCityGMLStatus = runningTag
	}
	if plateau {
		item.MergePlateauStatus = runningTag
	}
	if maxlod {
		item.MergeMaxLODStatus = runningTag
	}

	if err := c.UpdateDataItem(ctx, item); err != nil {
		log.Errorfc(ctx, "failed to update data item %s: %v", c.DataItemID, err)
	}
}

func (c *CMSWrapper) NotifyError(ctx context.Context, err error, citygml, plateau, maxlod bool) {
	if c == nil || err == nil {
		return
	}

	c.NotifyErrorMessage(ctx, err.Error(), citygml, plateau, maxlod)
}

func (c *CMSWrapper) NotifyErrorMessage(ctx context.Context, msg string, citygml, plateau, maxlod bool) {
	if c == nil {
		return
	}

	c.Comment(ctx, "公開準備処理に失敗しました。"+msg)

	item := &GspatialjpDataItem{}
	if citygml {
		item.MergeCityGMLStatus = failedTag
	} else if !c.SkipCityGML {
		item.MergeCityGMLStatus = idleTag
	}
	if plateau {
		item.MergePlateauStatus = failedTag
	} else if !c.SkipPlateau {
		item.MergePlateauStatus = idleTag
	}
	if maxlod {
		item.MergeMaxLODStatus = failedTag
	} else if !c.SkipMaxLOD {
		item.MergeMaxLODStatus = idleTag
	}

	if err := c.UpdateDataItem(ctx, item); err != nil {
		log.Errorfc(ctx, "failed to update data item %s: %v", c.DataItemID, err)
	}
}

func (c *CMSWrapper) UpdateDataItem(ctx context.Context, item *GspatialjpDataItem) error {
	if c == nil {
		return nil
	}

	p := &cms.Item{}
	cms.Marshal(item, p)
	p.ID = c.DataItemID

	if _, err := c.CMS.UpdateItem(ctx, p.ID, p.Fields, p.MetadataFields); err != nil {
		return err
	}

	return nil
}

func (c *CMSWrapper) UploadFile(ctx context.Context, path string) (string, error) {
	if c == nil {
		return "", nil
	}

	name := filepath.Base(path)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()
	return c.CMS.UploadAssetDirectly(ctx, c.ProjectID, name, f)
}

func (c *CMSWrapper) Upload(ctx context.Context, name string, body io.Reader) (string, error) {
	if c == nil {
		return "", nil
	}

	return c.CMS.UploadAssetDirectly(ctx, c.ProjectID, name, body)
}

func (c *CMSWrapper) Comment(ctx context.Context, comment string) {
	if c == nil {
		return
	}

	if err := c.CMS.CommentToItem(ctx, c.DataItemID, comment); err != nil {
		log.Errorfc(ctx, "failed to comment to %s: %v", c.DataItemID, err)
	}

	if err := c.CMS.CommentToItem(ctx, c.CityItemID, comment); err != nil {
		log.Errorfc(ctx, "failed to comment to %s: %v", c.CityItemID, err)
	}
}
