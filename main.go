package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type products struct {
	recettes
	Name string `json:"name"`
	URL  string `json:"url"`
}

type recettes struct {
	Descriptions string `json:"descriptions"`
	Ingredients  string `json:"ingredients"`
	Photos       string `json:"photos"`
	Directions   string `json:"directions"`
}

var allProducts []products
var allRecettes []recettes

func main() {
	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
	})

	// OnHTML enregistre une fonction. La fonction sera exécutée sur chaque HTML élément correspondant au paramètre
	c.OnHTML("a.mntl-card", func(h *colly.HTMLElement) {
		products := products{
			URL:  h.Attr("href"),
			Name: h.ChildText(".card__title-text"),
		}
		fmt.Println(products)
		allProducts = append(allProducts, products)
	})

	c.OnHTML("a.mntl-card", func(p *colly.HTMLElement) {
		nextPage := p.Request.AbsoluteURL(p.Attr("URL"))
		c.Visit(nextPage) // Je veux aller sur une autre page.
	})
 
	c.OnHTML("article.allrecipes-article", func(h *colly.HTMLElement) {
		recettes := recettes{
			Descriptions: h.ChildText("p.article-subheading"),
			Photos:       h.Attr("src"),
			Ingredients:  h.ChildText("li.mntl-structured-ingredients__list-item"),
			Directions:   h.ChildText("div.mntl-block"),
		}
		fmt.Println(recettes)
		allRecettes = append(allRecettes, recettes)
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
