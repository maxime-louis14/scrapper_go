package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type products struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Image string `json:"image"`
}

var allProducts []products

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.allrecipes.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
	})

	c.OnHTML("div.card__content", func(h *colly.HTMLElement) {
		products := products{
			Name: h.ChildText("span.card__title-text"),
			Image: h.ChildText("a[href]"),
		}
		fmt.Println(products)
		allProducts = append(allProducts, products)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "nError:", err)
	})

	c.Visit("https://www.allrecipes.com/")
}
