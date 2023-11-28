package cmd

import (
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
