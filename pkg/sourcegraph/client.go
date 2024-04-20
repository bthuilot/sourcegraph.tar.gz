package sourcegraph

import (
	"context"
	"fmt"
	"github.com/hasura/go-graphql-client"
	"net/http"
)

type Client interface {
	SearchFiles(query string) ([]FileSearchResult, error)
	GetFile(repository, path string, lines int) (File, error)
}

type client struct {
	gql *graphql.Client
	ctx context.Context
}

const GraphQLURL = "https://sourcegraph.com/.api/graphql"

type AdditionalHeaderTransport struct {
	T                 http.RoundTripper
	AdditionalHeaders map[string]string
}

func (adt *AdditionalHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range adt.AdditionalHeaders {
		req.Header.Add(k, v)
	}
	return adt.T.RoundTrip(req)
}

func NewClient(accessToken string) Client {
	httpClient := &http.Client{
		Transport: &AdditionalHeaderTransport{
			T: http.DefaultTransport,
			AdditionalHeaders: map[string]string{
				"Authorization": fmt.Sprintf("token %s", accessToken),
			},
		},
	}
	gql := graphql.NewClient(GraphQLURL, httpClient)
	return &client{
		gql: gql,
		ctx: context.Background(),
	}
}
