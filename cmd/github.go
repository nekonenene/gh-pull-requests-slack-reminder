package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

var (
	ctx          context.Context
	githubClient *github.Client
)

type ReviewState string

const (
	Approved         ReviewState = "APPROVED"
	ChangesRequested ReviewState = "CHANGES_REQUESTED"
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
// Returns a map of user IDs grouped by review state
func FetchReviewedUsersByIssue(issue *github.Issue) (map[string][]string, error) {
	var reviews []*github.PullRequestReview
	userAndReviewState := make(map[string]string)     // key: user ID, value: review state ("approved" or "changes_requested")
	usersEachReviewState := make(map[string][]string) // key: review state ("approved" or "changes_requested"), value: user IDs
	pageNum := 1

	for {
		tmpReviews, resp, err := githubClient.PullRequests.ListReviews(ctx, params.GitHubOwner, params.GitHubRepo, issue.GetNumber(), &github.ListOptions{
			PerPage: PerPageDefault,
			Page:    pageNum,
		})
		if err != nil {
			return usersEachReviewState, err
		}

		reviews = append(reviews, tmpReviews...)

		if resp.NextPage == 0 {
			break
		} else {
			pageNum = resp.NextPage
		}
	}

	// Order reviews by latest submitted_at
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].GetSubmittedAt().After(reviews[j].GetSubmittedAt().Time)
	})

	for _, review := range reviews {
		userId := review.User.GetLogin()

		if userAndReviewState[userId] != "" { // Skip if the user has already been added
			continue
		}

		if review.GetState() == "APPROVED" {
			userAndReviewState[userId] = "approved"
		}

		if review.GetState() == "CHANGES_REQUESTED" {
			userAndReviewState[userId] = "changes_requested"
		}
	}

	for userId, reviewState := range userAndReviewState {
		usersEachReviewState[reviewState] = append(usersEachReviewState[reviewState], userId)
	}

	// Reverse users to sort them by review submitted
	for _, users := range usersEachReviewState {
		sort.Slice(users, func(i, j int) bool {
			return users[i] < users[j]
		})
	}

	return usersEachReviewState, nil
}
