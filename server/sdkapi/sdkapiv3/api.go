package sdkapiv3

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hasura/go-graphql-client"
)

type GqlClient struct {
	client *graphql.Client
}

func NewClient(conf Config) (*GqlClient, error) {
	gqlURL, err := url.JoinPath(conf.BaseURL, "/datacatalog/graphql")
	if err != nil {
		return nil, fmt.Errorf("error joining base URL and graphql path: %w", err)
	}

	return &GqlClient{
		client: graphql.NewClient(gqlURL, nil),
	}, nil
}

func (c *GqlClient) QueryDatasets() (Query, error) {
	var q Query

	err := c.client.Query(context.Background(), &q, nil)
	if err != nil {
		return q, fmt.Errorf("error querying datasets: %w", err)
	}

	return q, nil
}
