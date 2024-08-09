package main

import (
	"fmt"
	"top40/input"
	"top40/scraper"
)

func main() {
	targetYear, targetMonth, targetWeekOfMonth := input.ParseYearMonthWeek()
	fmt.Printf("targetYear: %d, targetMonth: %d, targetWeekOfMonth: %d\n", targetYear, targetMonth, targetWeekOfMonth)

	errorBackNumber := scraper.BackNumber(targetYear, targetMonth, targetWeekOfMonth)
	if errorBackNumber != nil {
		fmt.Println(errorBackNumber)
		errorArchives := scraper.Archives(targetYear, targetMonth, targetWeekOfMonth)
		if errorArchives != nil {
			fmt.Println(errorArchives)
		}
	}
}
