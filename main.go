package main

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
)

type VisitedLinks struct {
	Url   map[string]bool
	Total int
}

type Collection struct {
	TotalCount int
	Data       []CollectedData
}

type CollectedData struct {
	Url    string            `json:"links"`
	Title  string            `json:"title"`
	Sku    string            `json:"sku"`
	Price  string            `json:"price"`
	Params map[string]string `json:"params"`
}

const (
	httpPart         = "http"
	maxNumberOfPages = 5
)

func main() {
	visited := VisitedLinks{}
	visited.addToVisit("https://example.com/", false)
	visited.parsePagesContent()
}

// Go through all pages marked for visiting and parse its content
// Mark page as visited afterward
func (v *VisitedLinks) parsePagesContent() {
	if !v.hasUnvisitedLinks() || v.Total > maxNumberOfPages {
		fmt.Println("There is no unvisited links.")
		return
	}
	c := colly.NewCollector(colly.AllowedDomains("example.com"))
	collection := Collection{
		TotalCount: 0,
	}

	for link, isVisited := range v.Url {
		if isVisited {
			continue
		}
		err, cd := v.VisitPage(c, link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		collection.TotalCount++
		collection.Data = append(collection.Data, *cd)
	}
	v.parsePagesContent()

	for _, data := range collection.Data {
		fmt.Println(
			data.Url,
			data.Title,
			data.Sku,
			data.Price,
			data.Params,
		)
	}
}

func (v *VisitedLinks) VisitPage(c *colly.Collector, link string) (error, *CollectedData) {
	cd := &CollectedData{}
	cd.Url = link

	shouldSkipPage := true
	//check if page is product page
	c.OnHTML("meta[property='product:product_link']", func(e *colly.HTMLElement) {
		shouldSkipPage = false
	})

	// Collect other links on page to be visited
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !checkUrlCanBeVisited(e.Attr("href")) {
			v.addToVisit(e.Attr("href"), false)
		}
	})

	// Find and save product name from page
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		cd.Title = e.Text
	})

	// find and save price
	c.OnHTML("div.price-action-group div[class='price'] span", func(e *colly.HTMLElement) {
		for _, node := range e.DOM.Nodes {
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "price-new" || attr.Val == "autocalc-product-price" {
					cd.Price = e.Text
				}
			}
		}
	})

	// find and save SKU (article)
	c.OnHTML("div.info-sku span", func(e *colly.HTMLElement) {
		cd.Sku = e.Text
	})

	params := make(map[string]string)
	//find and save characteristics info
	c.OnHTML("div#tab-specification .short-attribute", func(e *colly.HTMLElement) {
		params[e.ChildText("span.attr-name")] = e.ChildText("span.attr-text")
		cd.Params = params
	})

	err := c.Visit(link)
	v.markVisited(link)
	if err != nil || shouldSkipPage {
		return errors.New("error on visiting link"), nil
	}

	return nil, cd
}

// Check if we should skip url from adding to links pool
func checkUrlCanBeVisited(url string) bool {
	if url == "" || !strings.Contains(url, httpPart) {
		return true
	}
	skip := false
	for _, part := range getUrlPartsToSkip() {
		if strings.Contains(url, part) {
			skip = true
			break
		}
	}

	return skip
}

// Add link to pool of links with visited flag
func (v *VisitedLinks) addToVisit(url string, visited bool) {
	if v.Url == nil {
		v.Url = make(map[string]bool)
	}
	if !v.Url[url] {
		v.Url[url] = visited
	}
}

// function check if there is at least 1 link in pool that was not visited
func (v *VisitedLinks) hasUnvisitedLinks() bool {
	if v.Url == nil {
		return false
	}

	for _, isVisited := range v.Url {
		if isVisited == false {
			return true
		}
	}

	return false
}

// Mark given url as visited, so it will not be parsed again
func (v *VisitedLinks) markVisited(url string) {
	if v.Url == nil {
		v.Url = make(map[string]bool)
		return
	}

	if v.Url[url] && v.Url[url] == true {
		// already visited
		return
	}
	v.Total++
	v.Url[url] = true
}

// if url contains any of this parts - it will be skipped during links collection
func getUrlPartsToSkip() []string {
	urlPartsToSkip := []string{
		"cache",
		"image",
		"t.me",
		"wa.me",
	}

	return urlPartsToSkip
}
