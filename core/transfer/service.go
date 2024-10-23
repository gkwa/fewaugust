package transfer

import (
	"context"
	"fmt"

	"github.com/gkwa/fewaugust/core/repository"
	"github.com/go-logr/logr"
)

type Service interface {
	TransferPublicRepositories(ctx context.Context, currentOwner, newOwner string, repos, exclude []string, dryRun bool) error
}

type service struct {
	client repository.Client
	logger logr.Logger
}

func NewService(client repository.Client, logger logr.Logger) Service {
	return &service{
		client: client,
		logger: logger,
	}
}

func (s *service) TransferPublicRepositories(ctx context.Context, currentOwner, newOwner string, targetRepos, excludeRepos []string, dryRun bool) error {
	repos, err := s.client.ListRepositories(ctx, currentOwner)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if repo.Visibility != "public" {
			s.logger.V(1).Info("skipping non-public repository", "name", repo.Name)
			continue
		}

		// Check if repo is in exclude list
		excluded := false
		for _, excludeRepo := range excludeRepos {
			if repo.Name == excludeRepo {
				s.logger.V(1).Info("skipping excluded repository", "name", repo.Name)
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		// If specific repos were requested, check if this repo is in the list
		if len(targetRepos) > 0 {
			found := false
			for _, targetRepo := range targetRepos {
				if repo.Name == targetRepo {
					found = true
					break
				}
			}
			if !found {
				s.logger.V(1).Info("skipping repository not in target list", "name", repo.Name)
				continue
			}
		}

		currentURL := fmt.Sprintf("https://github.com/%s/%s", currentOwner, repo.Name)
		newURL := fmt.Sprintf("https://github.com/%s/%s", newOwner, repo.Name)

		if dryRun {
			s.logger.Info("[DRY RUN] would transfer repository",
				"name", repo.Name,
				"from", currentOwner,
				"to", newOwner,
				"current_url", currentURL,
				"new_url", newURL,
			)
			continue
		}

		s.logger.Info("transferring repository",
			"name", repo.Name,
			"from", currentOwner,
			"to", newOwner,
			"current_url", currentURL,
			"new_url", newURL,
		)

		if err := s.client.TransferRepository(ctx, repo, newOwner); err != nil {
			s.logger.Error(err, "failed to transfer repository", "name", repo.Name)
			continue
		}

		s.logger.Info("successfully transferred repository",
			"name", repo.Name,
			"current_url", newURL,
		)
	}

	return nil
}
