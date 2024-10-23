package repository

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Client interface {
	ListRepositories(ctx context.Context, owner string) ([]Repository, error)
	TransferRepository(ctx context.Context, repo Repository, newOwner string) error
}

type ghClient struct {
	client *api.RESTClient
}

func NewClient() (Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}
	return &ghClient{client: client}, nil
}
