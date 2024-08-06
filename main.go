package main

import (
	"fmt"
	"top40/input"
	"top40/scraper"
)

func main() {
	year, month, weekOfMonth := input.ParseYearMonthWeek()
	fmt.Printf("argsYear: %d, argsMonth: %d, argsWeek: %d\n", year, month, weekOfMonth)

	scraper.Archives()
}
