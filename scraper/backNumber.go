package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func BackNumber(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		e.ForEach("div.oa_list", func(_ int, d *colly.HTMLElement) {
			dataHtmlAttr := d.Attr("data-html")

			apiURL := "https://www.tvk-yokohama.com/top40/" + dataHtmlAttr

			parsedYear := ParseYear(apiURL)
			parsedDate := ParseDateBackNumber(apiURL)
			fmt.Println(parsedYear, parsedDate)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/backnumber.html")
}
