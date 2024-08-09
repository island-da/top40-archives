package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

func BackNumber(targetYear int, targetMonth int, targetWeekOfMonth int) {
	c := colly.NewCollector()

	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		var urls []string
		e.ForEach("div.oa_list", func(_ int, d *colly.HTMLElement) {
			dataHtmlAttr := d.Attr("data-html")
			urls = append(urls, "https://www.tvk-yokohama.com/top40/"+dataHtmlAttr)
		})
		match(urls, targetYear, targetMonth, targetWeekOfMonth)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.tvk-yokohama.com/top40/backnumber.html")
}

func match(urls []string, targetYear int, targetMonth int, targetWeekOfMonth int) {
	count := 1
	for i := len(urls) - 1; i >= 0; i-- {
		url := urls[i]
		parsedYear := ParseYear(url)
		parsedDate := ParseDateBackNumber(url)
		if parsedYear == targetYear && parsedDate == targetMonth {
			if count == targetWeekOfMonth {
				popUp(url)
				break
			}
			count++
		}
	}
}

func popUp(url string) {
	c := colly.NewCollector()

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		var title, artist string
		e.ForEach("td", func(i int, d *colly.HTMLElement) {
			if i == 2 {
				title = d.Text
			}
			if i == 3 {
				artist = d.Text
			}
		})
		fmt.Printf("%s / %s\n", artist, title)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("PopUp", r.URL)
	})

	c.Visit(url)
}
