package cmd

import (
	"fmt"
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
}

func isWeekend() bool {
	currentTime := time.Now()

	return currentTime.Weekday() == time.Saturday || currentTime.Weekday() == time.Sunday
}

func isJapaneseHoliday() bool {
	return holiday.IsHoliday(time.Now());
}
