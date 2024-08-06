package scraper

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func Archives(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	c.OnHTML("div.month", func(e *colly.HTMLElement) {
		count := 1
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

				var parsedYear, parsedDate int
				segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
				if len(segments) >= 2 {
					parsedYear, _ = strconv.Atoi(segments[len(segments)-2])
					dateWithExt := segments[len(segments)-1]
					parsedDate, _ = strconv.Atoi(strings.TrimSuffix(dateWithExt, path.Ext(dateWithExt))[:2])
				} else {
					fmt.Println("Not enough segments in path:", parsedURL.Path)
					return
				}

				if parsedYear == targetYear && parsedDate == targetMonth && count == targetWeekOfMonth {
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
							fmt.Printf("%s / %s\n", artistName, title)
						} else {
							fmt.Println("Not enough fields")
						}
					}
				}

			}
			count++
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/2022/archives.html")
}
