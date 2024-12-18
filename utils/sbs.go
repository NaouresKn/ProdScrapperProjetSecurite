package utils

import (
	"errors"
	"fmt"
	"log"
	"scrapper/schema"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapperFromSBS(url string, search string) ([]schema.ProductDetail, error) {
	if url == "" || search == "" {
		return nil, errors.New("invalid parameters")
	}

	c := colly.NewCollector()

	var products []schema.ProductDetail

	urlMatch := fmt.Sprintf("%s%s", url, search)

	c.OnHTML(".product-miniature", func(e *colly.HTMLElement) {
		name := e.ChildText("h6[itemprop='name']")
		price := e.ChildText("span.price")
		image := e.ChildAttr("img.tvproduct-hover-img", "src")
		link := e.ChildAttr("a.thumbnail.product-thumbnail", "href")

		price = strings.ReplaceAll(price, " ", "")
		price = strings.ReplaceAll(price, "TND", "")

		price = strings.Split(price, "\u00A0")[0]

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
