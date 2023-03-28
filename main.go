package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type products struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

	// OnHTML enregistre une fonction. La fonction sera exécutée sur chaque HTML élément correspondant au paramètre
	c.OnHTML("a.mntl-card", func(h *colly.HTMLElement) {
		products := products{
			URL:  h.ChildAttr("a.mntl-card-list-items", "href"),
			Name: h.ChildText(".card__title-text"),
		}
		fmt.Println(products)
		allProducts = append(allProducts, products)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "nError:", err)
	})

	c.Visit("https://www.allrecipes.com/recipes/17562/dinner/")

	content, err := json.Marshal(allProducts)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("data.json", content, 0644)
	fmt.Println("Total produts: ", len(allProducts))
}
