package cmd

import (
	"fmt"
	"log"
	"time"

	holiday "github.com/holiday-jp/holiday_jp-go"
)

func Exec() {
	ParseParameters()

	if params.AvoidWeekend && isWeekend() {
		fmt.Printf("\033[33mSkipped because today is a weekend.\033[0m\n")
		return
	}

	if params.AvoidJapaneseHolidays && isJapaneseHoliday() {
		fmt.Printf("\033[33mSkipped because today is a Japanese holiday.\033[0m\n")
		return
	}

	err := InitContextAndGitHubClient()
	if err != nil {
		log.Fatal(err)
	}

	issues, err := FetchLabelRelatedPullRequestIssues(FetchIssuesLimit)
	if err != nil {
		log.Fatal(err)
	}

	blocks, err := ConstructBlocksByIssues(issues)
	if err != nil {
		log.Fatal(err)
	}

	err = SendSlackMessage(blocks)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\033[32mProcessing succeeded!\033[0m")
}

// Check whether today is a weekend
func isWeekend() bool {
	currentTime := time.Now()

	return currentTime.Weekday() == time.Saturday || currentTime.Weekday() == time.Sunday
}

// Check whether today is a Japanese holiday
func isJapaneseHoliday() bool {
	return holiday.IsHoliday(time.Now())
}
