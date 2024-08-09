package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func BackNumber(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		stop := false
		e.ForEach("div.oa_list", func(_ int, d *colly.HTMLElement) {
			if stop {
				return
			}

			dataHtmlAttr := d.Attr("data-html")
			apiURL := "https://www.tvk-yokohama.com/top40/" + dataHtmlAttr

			parsedYear := ParseYear(apiURL)
			parsedDate := ParseDateBackNumber(apiURL)

			count := 1
			if parsedYear == targetYear && parsedDate == targetMonth {
				if count == targetWeekOfMonth {
					stop = true
					// TODO
				}
				count++
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/backnumber.html")
}
