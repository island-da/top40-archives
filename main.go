package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {

	currentDate := time.Now()
	lastMonthDate := currentDate.AddDate(0, -1, 0)
	lastYear, lastMonth, _ := lastMonthDate.Date()

	argsYear := flag.Int("year", lastYear, "year to scrape")
	argsMonth := flag.Int("month", int(lastMonth), "month to scrape")
	argsWeekOfMonth := flag.Int("week", 1, "week of month to scrape")

	flag.Parse()
	fmt.Printf("argsYear: %d, argsMonth: %d, argsWeek: %d\n", *argsYear, *argsMonth, *argsWeekOfMonth)

	c := colly.NewCollector()

	c.OnHTML("div.month", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, a *colly.HTMLElement) {
			onClickAttr := a.Attr("onclick")

			if strings.Contains(onClickAttr, "loadDataFile") {
				start := strings.Index(onClickAttr, "'") + 1
				end := strings.LastIndex(onClickAttr, "'")
				apiURL := "https://www.tvk-yokohama.com/top40/2022/" + onClickAttr[start:end]

				parsedURL, err := url.Parse(apiURL)
				if err != nil {
					log.Printf("Failed to parse URL %s: %v", apiURL, err)
					return
				}

				year, date := "", ""
				segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
				if len(segments) >= 2 {
					year = segments[len(segments)-2]
					dateWithExt := segments[len(segments)-1]
					date = strings.TrimSuffix(dateWithExt, path.Ext(dateWithExt))
				} else {
					fmt.Println("Not enough segments in path:", parsedURL.Path)
					return
				}

				res, err := http.Get(apiURL)
				if err != nil {
					fmt.Println("Failed to make API request:", err)
					return
				}
				defer res.Body.Close()
				time.Sleep(1 * time.Second)

				reader := csv.NewReader(res.Body)
				for {
					record, err := reader.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatalf("Failed to read CSV: %v", err)
					}

					joinedRecord := strings.Join(record, ",")
					fields := strings.Split(joinedRecord, ",")
					if len(fields) >= 2 {
						title := html.UnescapeString(fields[len(fields)-2])
						artistName := html.UnescapeString(fields[len(fields)-1])
						fmt.Printf("%s_%s %s / %s\n", year, date, artistName, title)
					} else {
						fmt.Println("Not enough fields")
					}
				}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/2022/archives.html")
}
