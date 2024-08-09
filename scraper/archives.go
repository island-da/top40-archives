package scraper

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func Archives(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	found := false
	c.OnHTML("div.month", func(e *colly.HTMLElement) {
		if found {
			return
		}

		count := 1
		e.ForEachWithBreak("a", func(_ int, a *colly.HTMLElement) bool {
			onClickAttr := a.Attr("onclick")

			if strings.Contains(onClickAttr, "loadDataFile") {
				start := strings.Index(onClickAttr, "'") + 1
				end := strings.LastIndex(onClickAttr, "'")
				apiURL := "https://www.tvk-yokohama.com/top40/2022/" + onClickAttr[start:end]

				parsedYear := ParseYear(apiURL)
				parsedDate := ParseDateArchives(apiURL)

				if parsedYear == targetYear && parsedDate == targetMonth && count == targetWeekOfMonth {
					found = true
					res, err := http.Get(apiURL)
					if err != nil {
						log.Fatalln("Failed to make API request:", err)
						return false
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
							log.Fatalln("Not enough fields")
						}
					}
				}
			}
			count++
			return true
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/2022/archives.html")
}
