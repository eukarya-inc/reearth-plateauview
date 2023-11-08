package cmsintegrationv3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type Config = cmsintegrationcommon.Config

const HandlerPath = "/notify_fme/v3"

func resultURL(conf *Config) string {
	return fmt.Sprintf("%s%s", conf.Host, HandlerPath)
}

type Services struct {
	FME  fmeInterface
	CMS  cms.Interface
	HTTP *http.Client
}

func NewServices(c Config) (s *Services, _ error) {
	s = &Services{}

	if !c.FMEMock {
		fmeURL := c.FMEURLV3
		if fmeURL == "" {
			return nil, errors.New("FME URL is not set")
		}

		resultURL, err := url.JoinPath(c.Host, "/notify_fme")
		if err != nil {
			return nil, fmt.Errorf("failed to init fme: %w", err)
		}

		fme := newFME(fmeURL, resultURL, c.FMESkipQualityCheck)
		s.FME = fme
	}

	cms, err := cms.New(c.CMSBaseURL, c.CMSToken)
	if err != nil {
		return nil, fmt.Errorf("failed to init cms: %w", err)
	}
	s.CMS = cms

	return
}

func (s *Services) UpdateFeatureItemStatus(ctx context.Context, itemID string, status ConvertionStatus) error {
	fields := (&FeatureItem{
		ConvertionStatus: status,
	}).CMSItem().MetadataFields
	_, err := s.CMS.UpdateItem(ctx, itemID, nil, fields)
	if err != nil {
		j, _ := json.Marshal(fields)
		log.Debugfc(ctx, "cmsintegrationv3: item update for %s: %s", itemID, j)
	}
	return err
}

func (s *Services) DownloadAsset(ctx context.Context, assetID string) (io.ReadCloser, error) {
	asset, err := s.CMS.Asset(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := s.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		_ = res.Body.Close()
		return nil, fmt.Errorf("failed to download asset: %s", res.Status)
	}

	return res.Body, nil
}

func (s *Services) DownloadAssetAsBytes(ctx context.Context, assetID string) ([]byte, error) {
	body, err := s.DownloadAsset(ctx, assetID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = body.Close()
	}()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset: %w", err)
	}

	return buf.Bytes(), nil
}
