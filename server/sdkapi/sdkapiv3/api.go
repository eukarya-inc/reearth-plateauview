package sdkapiv3

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type GqlClient struct {
	client *graphql.Client
}

func NewClient(conf Config) (*GqlClient, error) {
	c := graphql.NewClient(conf.GQLBaseURL, nil).WithRequestModifier(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+conf.GQLToken)
	})

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

	vars := map[string]any{
		"code": AreaCode(id),
	}

	err := c.client.Query(context.Background(), &q, vars)
	if err != nil {
		return q, fmt.Errorf("error querying dataset files: %w", err)
	}

	return q, nil
}
