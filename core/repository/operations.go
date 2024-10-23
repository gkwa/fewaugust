package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

func (c *ghClient) ListRepositories(ctx context.Context, owner string) ([]Repository, error) {
	page := 1
	perPage := 100
	var allRepos []Repository

	for {
		var repos []Repository
		path := fmt.Sprintf("users/%s/repos?page=%d&per_page=%d", owner, page, perPage)
		resp, err := c.client.Request("GET", path, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		page++
	}

	return allRepos, nil
}

func (c *ghClient) TransferRepository(ctx context.Context, repo Repository, newOwner string) error {
	path := fmt.Sprintf("repos/%s/%s/transfer", repo.Owner.Login, repo.Name)
	payload := struct {
		NewOwner string `json:"new_owner"`
	}{
		NewOwner: newOwner,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.client.Post(path, bytes.NewReader(body), nil)
}
