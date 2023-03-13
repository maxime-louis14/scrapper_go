package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Product struct {
	Titre       string
	Image       string
	Description string
	Url         string
	Ingredient  string
	Preparation string
}

func main() {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	products := make([]Product, 0)

	// Callbacks
	c.OnHTML("a.[comp card--image-top mntl-card-list-items mntl-document-card mntl-card card card--no-image]", func(e *colly.HTMLElement) {
		e.ForEach("span.[card__title-text]", func(i int, h *colly.HTMLElement) {
			item := Product{}
			item.Titre = h.Text
			item.Image = e.ChildAttr("img", "data-src")
			item.Description = e.Attr("description")
			item.Url = "https://www.allrecipes.com/" + e.Attr(("href"))
			item.Preparation = e.ChildText("div.tag._dsct")
			products = append(products, item)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Code", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finshed", r.Request.URL)
		js, err := json.MarshalIndent(products, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("products.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}
	})

	c.Visit("https://www.allrecipes.com/")
}
