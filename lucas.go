package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"github.com/gocolly/colly"
)

type Clothing struct {
	Name					string
	Code					string
	Description		string
	Price					string
}


func main() {
	c := colly.NewCollector(
		// TODO modify with regex later for optimized crawling...
		// colly.AllowedDomains("https://www.floryday.com/"),
		colly.CacheDir(".floryday_cache"),
  	colly.MaxDepth(3), // keeping crawling limited for our initial experiments
  )

	// clothing detail scraping collector
	detailCollector := c.Clone()

	clothes := make([]Clothing, 0, 200)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		// debug log
		log.Print("LINK TO HIT", link)

		// hardcoded urls to skip -> to be optimized -> perhaps map links from external file...
		if !strings.HasPrefix(link, "/?country_code") || strings.Index(link, "/cart.php") > -1 ||
		strings.Index(link, "/login.php") > -1 || strings.Index(link, "/cart.php") > -1 ||
		strings.Index(link, "/account") > -1 || strings.Index(link, "/privacy-policy.html") > -1 {
			return
		}

		// scrape the page
		e.Request.Visit(link)
	})

	// printing visiting message for debug purposes
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String(), "\n")
	})

	// TODO filter this better a[href] is way too broad -> may need regex
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {

		clothingURL := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Validating crawling Link: ", clothingURL)

		// links provided need to be better filtered
		// hardcoding one value only to work here for now...
		if clothingURL == "https://www.floryday.com/Cotton-Floral-Short-Sleeve-High-Low-Dress-m1043239" {
			// Activate detailCollector
			log.Println("Crawling Link Validated -> Commencing Crawl...")
			detailCollector.Visit(clothingURL)
		} else {
			log.Println("Validation Failed -> Cancelling Crawl...")
			return
		}

	})

	// Extract details of the clothing
	detailCollector.OnHTML(`div[class=prod-right-in]`, func(e *colly.HTMLElement) {

		// TODO secure variables with default error strings in event values are missing
		title := e.ChildText(".prod-name")
		code := e.ChildText(".prod-item-code")
		price := e.ChildText(".prod-price")
		description := e.ChildText(".grid-uniform")

		clothing := Clothing{
			Name: 					title,
			Code: 					code,
			Description: 		description,
			Price:					price,
		}

		clothes = append(clothes, clothing)
	})

	// start scraping at our seed address
	c.Visit("https://www.floryday.com/Dresses-r9872/")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	// TODO this could be dumped into a DB instead
	enc.Encode(clothes)

}
