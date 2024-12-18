package utils

import (
	"errors"
	"fmt"
	"log"
	"scrapper/schema"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapperFromTunisianet(url string, search string) ([]schema.ProductDetail, error) {
	if url == "" || search == "" {
		return nil, errors.New("invalid parameters")
	}

	c := colly.NewCollector()

	var products []schema.ProductDetail

	urlMatch := fmt.Sprintf("%s%s", url, search)

	c.OnHTML(".product-miniature", func(e *colly.HTMLElement) {

		name := e.ChildText("h2.h3.product-title a")
		link := e.ChildAttr("h2.h3.product-title a", "href")
		image := e.ChildAttr("img.center-block.img-responsive", "src")
		price := e.ChildText("span.price")

		price = strings.TrimSpace(price)
		price = strings.ReplaceAll(price, "DT", "")

		if name != "" && price != "" && image != "" && link != "" {
			product := schema.ProductDetail{
				Name:  strings.TrimSpace(name),
				Price: strings.TrimSpace(price),
				Image: strings.TrimSpace(image),
				Link:  strings.TrimSpace(link),
			}
			products = append(products, product)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", err)
	})

	err := c.Visit(urlMatch)
	if err != nil {
		return nil, err
	}

	c.Wait()
	return products, nil
}