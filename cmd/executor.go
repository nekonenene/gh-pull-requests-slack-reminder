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

	if len(issues) == 0 {
		fmt.Printf("\033[33mSkipped because there are no Pull Requests related the specified label in the open state.\033[0m\n")
		return
	}

	fmt.Println("Issues count: ", len(issues))

	issuesEachAuthor := IssuesEachAuthor(issues)

	fmt.Println("Issue Each Author: ", issuesEachAuthor)

	fmt.Println("Processing succeeded!")
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
