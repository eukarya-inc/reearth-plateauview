package datacatalogv2adapter

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func fetchMaxLOD(ctx context.Context, all []datacatalogv2.DataCatalogItem) error {
	urls := lo.Map(all, func(item datacatalogv2.DataCatalogItem, _ int) string {
		return item.MaxLODURL
	})

	maxlod, err := fetchMaxLODContents(ctx, urls)
	if err != nil {
		return fmt.Errorf("failed to fetch max lod: %w", err)
	}

	for i, m := range maxlod {
		all[i].MaxLODContent = m
	}
	return nil
}

func fetchMaxLODContents(ctx context.Context, urls []string) ([][][]string, error) {
	res := make([][][]string, len(urls))
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(10)

	for i := 0; i < len(urls); i++ {
		i := i
		url := urls[i]
		if url == "" {
			continue
		}

		eg.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("items[%d]: failed to create request: %w", i, err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("items[%d]: failed to get max LOD content: %w", i, err)
			}

			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("items[%d]: failed to get max LOD content: status code %d", i, resp.StatusCode)
			}

			c := csv.NewReader(resp.Body)
			records, err := c.ReadAll()
			if err != nil {
				return fmt.Errorf("items[%d]: failed to read max LOD content: %w", i, err)
			}

			res[i] = records
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}
