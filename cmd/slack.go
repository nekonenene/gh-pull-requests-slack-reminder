package cmd

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v56/github"
	"github.com/slack-go/slack"
)

// Send a message to Slack using Incoming Webhook
func SendSlackMessage(blocks *slack.Blocks) error {
	err := slack.PostWebhook(params.SlackWebhookUrl, &slack.WebhookMessage{
		Blocks: blocks,
	})
	if err != nil {
		return err
	}

	return nil
}

func ConstructBlocksByIssues(issues []*github.Issue) (*slack.Blocks, error) {
	if len(issues) == 0 {
		blocks := &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", "*No pull requests to review!!* 👏", true, false),
					nil,
					nil,
				),
			},
		}

		return blocks, nil
	}

	blocks := &slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "*Please review the following PRs!* 😎", true, false),
				nil,
				nil,
			),
			slack.NewDividerBlock(),
		},
	}

	issuesEachAuthor := IssuesEachAuthor(issues)

	for authorUserId, issues := range issuesEachAuthor {
		checkboxesObjects := make([]*slack.OptionBlockObject, len(issues))

		// Construct checkboxes
		for i, issue := range issues {
			approvedUsers, err := FetchApprovedUsersByIssue(issue)
			if err != nil {
				return nil, err
			}

			text := fmt.Sprintf("*<%s|#%d %s>*", issue.GetURL(), issue.GetNumber(), issue.GetTitle())

			var description *slack.TextBlockObject
			if len(approvedUsers) > 0 {
				description = slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Approved by %s", strings.Join(approvedUsers, ", ")), false, false)
			}

			checkboxesObjects[i] = &slack.OptionBlockObject{
				Text:        slack.NewTextBlockObject("mrkdwn", text, false, false),
				Description: description,
				Value:       fmt.Sprintf("pr-%v", issue.GetNumber()),
			}
		}

		pullRequestsBlocks := &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", authorUserId), false, false),
					nil,
					nil,
				),
				slack.NewActionBlock(
					fmt.Sprintf("%s-checkboxes-action", authorUserId),
					slack.NewCheckboxGroupsBlockElement(
						fmt.Sprintf("%s-checkboxes", authorUserId),
						checkboxesObjects...,
					),
				),
			},
		}

		jsonBlocks, _ := pullRequestsBlocks.MarshalJSON()
		fmt.Println("Pull Request Blocks: ", string(jsonBlocks))

		blocks.BlockSet = append(blocks.BlockSet, pullRequestsBlocks.BlockSet...)
	}

	jsonBlocks, _ := blocks.MarshalJSON()
	fmt.Println("Blocks: ", string(jsonBlocks))

	return blocks, nil
}
