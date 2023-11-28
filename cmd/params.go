package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type parameters struct {
	GitHubAPIToken        string
	GitHubOwner           string
	GitHubRepo            string
	GitHubEnterpriseURL   string
	TargetLabelName       string
	SlackWebhookUrl       string
	AvoidWeekend          bool
	AvoidJapaneseHolidays bool
	ShowVersion           bool
}

var params parameters
var Version = "" // Overwrite when building

func ParseParameters() {
	flag.StringVar(&params.GitHubAPIToken, "token", "", "[Required] GitHub API Token")
	flag.StringVar(&params.GitHubOwner, "owner", "", "[Required] Owner name of the repository (e.g. octocat)")
	flag.StringVar(&params.GitHubRepo, "repo", "", "[Required] Repository name (e.g. hello-world)")
	flag.StringVar(&params.GitHubEnterpriseURL, "enterprise-url", "", "[Opiton] URL of GitHub Enterprise (ex. https://github.your.domain )")
	flag.StringVar(&params.TargetLabelName, "label-name", "", "[Required] Label name related to target pull requests (e.g. in-review)")
	flag.StringVar(&params.SlackWebhookUrl, "webhook-url", "", "[Required] URL of Slack Incoming Webhook (e.g. https://hooks.slack.com/services/XXXXXX)")
	flag.BoolVar(&params.AvoidWeekend, "avoid-weekend", false, "[Opiton] If true, don't send notifications on weekends")
	flag.BoolVar(&params.AvoidJapaneseHolidays, "avoid-jp-holidays", false, "[Opiton] If true, don't send notifications on Japanese holidays")
	flag.BoolVar(&params.ShowVersion, "version", false, "[Opiton] Show version")
	flag.BoolVar(&params.ShowVersion, "v", false, "[Opiton] Shorthand of -version")
	flag.Parse()

	if params.ShowVersion {
		fmt.Println(getVersionString())
		os.Exit(0)
	}

	// Validations
	if params.GitHubAPIToken == "" {
		log.Fatalln("-token is required")
	}
	if params.GitHubOwner == "" {
		log.Fatalln("-owner is required")
	}
	if params.GitHubRepo == "" {
		log.Fatalln("-repo is required")
	}
	if params.TargetLabelName == "" {
		log.Fatalln("-label-name is required")
	}
	if params.SlackWebhookUrl == "" {
		log.Fatalln("-webhook-url is required")
	}
}

func getVersionString() string {
	// For downloading a binary from GitHub Releases
	if Version != "" {
		return Version
	}

	// For `go install`
	i, ok := debug.ReadBuildInfo()
	if ok {
		return i.Main.Version
	}

	return "Development version"
}
