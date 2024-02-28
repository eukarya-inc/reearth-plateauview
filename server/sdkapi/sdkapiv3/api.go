package sdkapiv3

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hasura/go-graphql-client"
	"github.com/reearth/reearthx/log"
)

type GqlClient struct {
	client *graphql.Client
}

func NewClient(conf Config) (*GqlClient, error) {
	gqlURL, err := url.JoinPath(conf.BaseURL, "/datacatalog/admin/graphql")
	if err != nil {
		return nil, fmt.Errorf("error joining base URL and graphql path: %w", err)
	}

	c := graphql.NewClient(gqlURL, nil).WithRequestModifier(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+conf.GQLToken)
	})

	log.Infof("NewClient: %v", c)

	return &GqlClient{
		client: c,
	}, nil
}

func (c *GqlClient) QueryDatasets() (DatasetsQuery, error) {
	var q DatasetsQuery

	err := c.client.Query(context.Background(), &q, nil)
	if err != nil {
		return q, fmt.Errorf("error querying datasets: %w", err)
	}

	return q, nil
}

func (c *GqlClient) QueryDatasetFiles(id string) (DatasetFilesQuery, error) {
	var q DatasetFilesQuery

	vars := map[string]interface{}{
		"code": graphql.String(id),
	}

	err := c.client.Query(context.Background(), &q, vars)
	if err != nil {
		return q, fmt.Errorf("error querying dataset files: %w", err)
	}

	return q, nil
}
