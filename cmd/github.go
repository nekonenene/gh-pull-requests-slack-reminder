package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

var (
	ctx          context.Context
	githubClient *github.Client
)

// Init ctx and githubClient
func InitContextAndGitHubClient() error {
	ctx = context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	if params.GitHubEnterpriseURL == "" {
		githubClient = github.NewClient(httpClient)
	} else {
		githubClient, err := github.NewEnterpriseClient(params.GitHubEnterpriseURL, params.GitHubEnterpriseURL, httpClient)
		if err != nil {
			return err
		}
	}

	return nil
}
