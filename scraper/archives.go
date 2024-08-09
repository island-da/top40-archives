package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func Archives(targetYear int, targetMonth int, targetWeekOfMonth int) error {
	c := colly.NewCollector()

	founded := false
	c.OnHTML("div.month", func(e *colly.HTMLElement) {
		var urls []string
		e.ForEach("a", func(_ int, a *colly.HTMLElement) {
			onClickAttr := a.Attr("onclick")
			if strings.Contains(onClickAttr, "loadDataFile") {
				start := strings.Index(onClickAttr, "'") + 1
				end := strings.LastIndex(onClickAttr, "'")
				url := "https://www.tvk-yokohama.com/top40/2022/" + onClickAttr[start:end]
				urls = append(urls, url)
			}
		})
		founded = founded || matchArchives(urls, targetYear, targetMonth, targetWeekOfMonth)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/2022/archives.html")

	if !founded {
		return errors.New("archives not found")
	}
	return nil
}

func matchArchives(urls []string, targetYear int, targetMonth int, targetWeekOfMonth int) bool {
	founded := false
	for count, url := range urls {
		parsedYear := ParseYear(url)
		parsedDate := ParseDateArchives(url)
		if parsedYear == targetYear && parsedDate == targetMonth && count+1 == targetWeekOfMonth {
			founded = true
			csvReader(url)
			break
		}
	}
	return founded
}

func csvReader(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("Failed to make API request:", err)
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
			artist := html.UnescapeString(fields[len(fields)-1])
			fmt.Printf("%s / %s\n", artist, title)
		} else {
			log.Fatalln("Not enough fields")
		}
	}
}
