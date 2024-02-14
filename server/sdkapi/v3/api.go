package sdkapiv3

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

type Client struct {
	client *graphql.Client
}

func NewClient(conf Config) *Client {
	return &Client{
		client: graphql.NewClient(conf.CMSBaseURL+"/datacatalog/graphql", nil),
	}
}

func (c Client) QueryDatasets() (Query, error) {
	var q Query

	err := c.client.Query(context.Background(), &q, nil)
	if err != nil {
		return q, fmt.Errorf("error querying datasets: %w", err)
	}

	return q, nil
}
