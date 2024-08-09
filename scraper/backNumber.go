package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func BackNumber(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	found := false
	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		if found {
			return
		}
		e.ForEachWithBreak("div.oa_list", func(_ int, d *colly.HTMLElement) bool {
			dataHtmlAttr := d.Attr("data-html")
			url := "https://www.tvk-yokohama.com/top40/" + dataHtmlAttr

			parsedYear := ParseYear(url)
			parsedDate := ParseDateBackNumber(url)

			count := 1
			if parsedYear == targetYear && parsedDate == targetMonth {
				if count == targetWeekOfMonth {
					found = true
					popUp(url)
					return false
				}
				count++
			}
			return true
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/backnumber.html")
}

func popUp(url string) {
	c := colly.NewCollector()

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		fmt.Println("h3 found:", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("PopUp", r.URL)
	})

	c.Visit(url)
}
