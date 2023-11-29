package cmd

import (
	"fmt"
	"strings"
	"time"

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
					slack.NewTextBlockObject("mrkdwn", "*No pull requests to review!!* 👏", false, false),
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
				slack.NewTextBlockObject("mrkdwn", "*Please review the following PRs!* 😎", false, false),
				nil,
				nil,
			),
			slack.NewDividerBlock(),
		},
	}

	issuesEachAuthor := IssuesEachAuthor(issues)

	for authorUserId, issues := range issuesEachAuthor {
		authorAndPullRequestsBlocks := &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(
					slack.NewTextBlockObject("plain_text", authorUserId, false, false),
				),
			},
		}

		pullRequestsBlocks := make([]*slack.RichTextBlock, len(issues))

		// Construct checkboxes
		for i, issue := range issues {
			pullRequestsBlocks[i] = slack.NewRichTextBlock(
				fmt.Sprintf("pr-%v-block-%d", issue.GetNumber(), time.Now().Unix()),
				slack.NewRichTextSection(
					slack.NewRichTextSectionTextElement(
						"▶️ ",
						nil,
					),
					slack.NewRichTextSectionLinkElement(
						issue.GetPullRequestLinks().GetHTMLURL(),
						fmt.Sprintf("#%d %s", issue.GetNumber(), issue.GetTitle()),
						&slack.RichTextSectionTextStyle{
							Bold: true,
						},
					),
				),
			)

			approvedUsers, err := FetchApprovedUsersByIssue(issue)
			if err != nil {
				return nil, err
			}

			if len(approvedUsers) > 0 {
				pullRequestsBlocks[i].Elements = append(pullRequestsBlocks[i].Elements, slack.NewRichTextSection(
					slack.NewRichTextSectionTextElement(
						fmt.Sprintf("(%s already approved)", strings.Join(approvedUsers, ", ")),
						&slack.RichTextSectionTextStyle{
							Italic: true,
						},
					),
				))
			}

			authorAndPullRequestsBlocks.BlockSet = append(authorAndPullRequestsBlocks.BlockSet, pullRequestsBlocks[i])
		}

		// jsonBlocks, _ := pullRequestsBlocks.MarshalJSON()
		// fmt.Println("Pull Request Blocks: ", string(jsonBlocks))

		blocks.BlockSet = append(blocks.BlockSet, authorAndPullRequestsBlocks.BlockSet...)
	}

	// jsonBlocks, _ := blocks.MarshalJSON()
	// fmt.Println("Blocks: ", string(jsonBlocks))

	return blocks, nil
}
