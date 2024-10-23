package cmd

import (
	"strings"

	"github.com/gkwa/fewaugust/core/repository"
	"github.com/gkwa/fewaugust/core/transfer"
	"github.com/spf13/cobra"
)

var (
	currentOwner string
	newOwner     string
	repoNames    string
	excludeRepos string
	dryRun       bool
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Transfer GitHub repositories between owners",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := LoggerFrom(cmd.Context())

		client, err := repository.NewClient()
		if err != nil {
			return err
		}

		service := transfer.NewService(client, logger)

		var repos []string
		if repoNames != "" {
			repos = strings.Split(repoNames, ",")
		}

		var exclude []string
		if excludeRepos != "" {
			exclude = strings.Split(excludeRepos, ",")
		}

		return service.TransferPublicRepositories(cmd.Context(), currentOwner, newOwner, repos, exclude, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVar(&currentOwner, "from", "", "Current repository owner")
	helloCmd.Flags().StringVar(&newOwner, "to", "", "New repository owner")
	helloCmd.Flags().StringVar(&repoNames, "repos", "", "Comma-separated list of repository names to transfer. If empty, all public repos will be transferred")
	helloCmd.Flags().StringVar(&excludeRepos, "exclude-repos", "", "Comma-separated list of repository names to exclude from transfer")
	helloCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions that would be taken without actually transferring repositories")
	if err := helloCmd.MarkFlagRequired("from"); err != nil {
		panic(err)
	}
	if err := helloCmd.MarkFlagRequired("to"); err != nil {
		panic(err)
	}
}
