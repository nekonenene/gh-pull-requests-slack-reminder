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
		var err error
		githubClient, err = github.NewClient(httpClient).WithEnterpriseURLs(params.GitHubEnterpriseURL, params.GitHubEnterpriseURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fetches all Pull Requests related the specified label in the open state as Issues.
// Only the number of issues defined by `limit` will be retrieved.
// NOTE: GitHub API allows only searching as Issues
func FetchLabelRelatedPullRequestIssues(limit int) ([]*github.Issue, error) {
	var issues []*github.Issue
	remainingLimit := limit
	pageNum := 1

	for {
		perPage := PerPageDefault
		if remainingLimit < PerPageDefault {
			perPage = remainingLimit
		}

		listOptions := github.ListOptions{
			PerPage: perPage,
			Page:    pageNum,
		}

		searchResult, resp, err := githubClient.Search.Issues(
			ctx,
			fmt.Sprintf("repo:%s/%s is:pr is:open label:\"%s\"", params.GitHubOwner, params.GitHubRepo, params.TargetLabelName),
			&github.SearchOptions{
				Sort:        "updated",
				Order:       "desc",
				ListOptions: listOptions,
			},
		)
		if err != nil {
			return issues, err
		}

		tmpIssues := searchResult.Issues
		issues = append(issues, tmpIssues...)

		// Stop fetching if the number of retrieved issues exceeds the limit
		remainingLimit -= len(tmpIssues)
		if remainingLimit == 0 {
			break
		}

		if resp.NextPage == 0 {
			break
		} else {
			pageNum = resp.NextPage
		}
	}

	return issues, nil
}

// Group issues by author
func IssuesEachAuthor(issues []*github.Issue) map[string][]*github.Issue {
	IssuesEachAuthor := make(map[string][]*github.Issue)

	for _, issue := range issues {
		userId := issue.User.GetLogin()
		IssuesEachAuthor[userId] = append(IssuesEachAuthor[userId], issue)
	}

	return IssuesEachAuthor
}

// Fetch user IDs who approved or requested changes the pull request
func FetchReviewedUsersByIssue(issue *github.Issue) (map[string][]string, error) {
	var reviews []*github.PullRequestReview
	reviewedUsers := make(map[string][]string) // key is "approved" or "changes_requested"
	pageNum := 1

	for {
		tmpReviews, resp, err := githubClient.PullRequests.ListReviews(ctx, params.GitHubOwner, params.GitHubRepo, issue.GetNumber(), &github.ListOptions{
			PerPage: PerPageDefault,
			Page:    pageNum,
		})
		if err != nil {
			return reviewedUsers, err
		}

		reviews = append(reviews, tmpReviews...)

		if resp.NextPage == 0 {
			break
		} else {
			pageNum = resp.NextPage
		}
	}

	for _, review := range reviews {
		if review.GetState() == "APPROVED" {
			reviewedUsers["approved"] = append(reviewedUsers["approved"], review.User.GetLogin())
		}

		if review.GetState() == "CHANGES_REQUESTED" {
			reviewedUsers["changes_requested"] = append(reviewedUsers["changes_requested"], review.User.GetLogin())
		}
	}

	return reviewedUsers, nil
}
