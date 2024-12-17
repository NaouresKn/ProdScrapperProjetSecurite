package scrapers

import (
	"fmt"

	"github.com/gocolly/colly"
)

func ScrapeSites(productName string) []string {
	var results []string

	sites := []string{
		fmt.Sprintf("https://example.com/search?q=%s", productName),
		fmt.Sprintf("https://another-example.com/search?q=%s", productName),
	}

	for _, site := range sites {
		c := colly.NewCollector()

		c.OnHTML("div.product", func(e *colly.HTMLElement) {
			product := e.ChildText("h2.name")
			price := e.ChildText("span.price")
			results = append(results, fmt.Sprintf("%s - %s", product, price))
		})

		// Visit each site
		c.Visit(site)
	}

	return results
}
