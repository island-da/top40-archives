package input

import (
	"flag"
	"time"
)

func ParseYearMonthWeek() (int, int, int) {
	currentDate := time.Now()
	lastMonthDate := currentDate.AddDate(0, -1, 0)
	lastYear, lastMonth, _ := lastMonthDate.Date()

	argsYear := flag.Int("year", lastYear, "year to scrape")
	argsMonth := flag.Int("month", int(lastMonth), "month to scrape")
	argsWeekOfMonth := flag.Int("week", 1, "week of month to scrape")

	flag.Parse()

	return *argsYear, *argsMonth, *argsWeekOfMonth
}
