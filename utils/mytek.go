package utils

import (
	"errors"
	"fmt"
	"log"
	"scrapper/schema"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapperFromMytek(url string, search string) ([]schema.ProductDetail, error) {
	if url == "" || search == "" {
		return nil, errors.New("invalid parameters")
	}

	c := colly.NewCollector()

	var products []schema.ProductDetail

	urlMatch := url + search

	c.OnHTML(".product-item-info", func(e *colly.HTMLElement) {
		name := e.ChildText(".product-item-link")
		price := e.ChildText(".price")
		image := e.ChildAttr(".product-image-photo", "src")
		link := e.ChildAttr(".product-item-link", "href")

		price = strings.ReplaceAll(price, " ", "")
		price = strings.ReplaceAll(price, "TND", "")

		price = strings.Split(price, " ")[0]

		if name != "" && price != "" && image != "" && link != "" {
			product := schema.ProductDetail{Name: name, Price: price, Image: image, Link: link}
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
