package main

import (
	"fmt"
	"top40/input"
	"top40/scraper"
)

func main() {
	targetYear, targetMonth, targetWeekOfMonth := input.ParseYearMonthWeek()
	fmt.Printf("targetYear: %d, targetMonth: %d, targetWeekOfMonth: %d\n", targetYear, targetMonth, targetWeekOfMonth)

	scraper.Archives(targetYear, targetMonth, targetWeekOfMonth)
}
