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

	c.OnHTML(".three-post--home", func(h *colly.HTMLElement) {
		products := products{
			Name: h.ChildText(".card__title-text"),
			URL: h.ChildAttr(".card--image-top", "href"),
		}
		fmt.Println(products)
		allProducts = append(allProducts, products)
	})

	// Ici le récupère les URL des pages
	// c.OnHTML(".three-post__inner", func(h *colly.HTMLElement) {
	// 	products := products{
	// 		URL: h.ChildAttr("a.card--image-top", "href"),
	// 	}
	// 	fmt.Println(products)
	// 	allProducts = append(allProducts, products)
	// })
	
	// je veux récupère les data d'une autre page
	// c.OnHTML("nav.header-nav", func(p *colly.HTMLElement) {
	// 	nextPage := p.Request.AbsoluteURL(p.Attr("li.header-nav__list-item"))
	// 	c.Visit(nextPage)
	// })

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "nError:", err)
	})

	c.Visit("https://www.allrecipes.com/")

	content, err := json.Marshal(allProducts)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("data.json", content, 0644)
	fmt.Println("Total produts: ", len(allProducts))
}
