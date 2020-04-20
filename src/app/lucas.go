package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
	"github.com/gocolly/colly"
	"github.com/fatih/color"
	"strconv"
	"github.com/joho/godotenv"
)

func main() {

	// loading config
	err := godotenv.Load()
  if err != nil {
    color.Red("Error loading .env file")
  }

	// setting up colly collector
	c := colly.NewCollector(
		// colly.AllowedDomains("https://www.floryday.com/"),
		colly.CacheDir(".floryday_cache"),
  	// colly.MaxDepth(5), // keeping crawling limited for our initial experiments,
		// colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
  )

	// clothing detail scraping collector
	detailCollector := c.Clone()

	// defaulting to array of 200 results
	clothes := make([]Clothing, 0, 200)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// hardcoded urls to skip -> to be optimized -> perhaps map links from external file...
		if !strings.HasPrefix(link, "/?country_code") || strings.Index(link, "/cart.php") > -1 ||
		strings.Index(link, "/login.php") > -1 || strings.Index(link, "/account") > -1 ||
		strings.Index(link, "/privacy-policy.html") > -1 {
			return
		}
		// scrape the page
		e.Request.Visit(link)
	})

	// printing visiting message for debug purposes
	c.OnRequest(func(r *colly.Request) {
		color.Blue("Visiting %s %s", r.URL.String(), "\n")
	})

	// TODO filter this better a[href] is way too broad -> may need regex
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		clothingURL := e.Request.AbsoluteURL(e.Attr("href"))
		// links provided need to be better filtered
		// hardcoding one value only to work here for now...
		if strings.Contains(clothingURL, "-Dress-"){
			// Activate detailCollector
			// Setting default country_code for currency purposes
			color.Magenta("Commencing Crawl for %s", clothingURL + "?country_code=IE")
			detailCollector.Visit(clothingURL + "?country_code=IE")
		} else {
			color.Red("Validation Failed -> Cancelling Crawl for %s", clothingURL + "?country_code=IE")
			return
		}
	})

	// experimental getting image url
	detailCollector.OnHTML(`.swipe-wrap`, func(e *colly.HTMLElement) {
		 link := e.ChildAttr("img", "src")
		 // ignore secondary blank source image element
		 if len(link) > 0 {
			 color.Blue("https:%s", e.ChildAttr("img", "src"))
		 }
	})

	// Extract details of the clothing (- image url)
	detailCollector.OnHTML(`html`, func(e *colly.HTMLElement) {

		// TODO secure variables with default error strings in event values are missing
		title := e.ChildText(".prod-name")
		code := strings.Split(e.ChildText(".prod-item-code"), "#")[1]

		// price parsing & reformatting
		initialprice := e.ChildText(".currency-prices")
		pricenosymbol := strings.TrimSuffix(initialprice," €")
		stringPrice := strings.Replace(pricenosymbol, ",", ".", 1)
		price, err := strconv.ParseFloat(stringPrice, 64) // conversion to float64
		if err != nil {
	    color.Red("err in parsing price -> %s", err)
	  }

		// desecription requires more refined parsing into subsections
		// description := strings.TrimSpace(e.ChildText(".grid-uniform"))

		url := "http://example.com"
		description := "Le Placeholder"
		photo := "random-image-url" // see logged result above and integrate logic

		clothing := Clothing{
			Name: title,
			Price: price,
			Url: url, // placeholder
			Description: description, // placeholder
			Code: code,
			Style: "style",
			Pattern: "pattern",
			Sleeve: "sleeve",
			Silhouette: "silhouette",
			Season: "season",
			Material: "material",
			Type: "type",
			Neckline: "neckline",
			Length: "length",
			Occasion: "occasion",
			Image: photo, // placeholder
			Date: time.Now().String(),
		}

		// writing as we go to DB
		// TODO optiize to handle bulk array uplaods instead
		dbWrite(clothing)

		// appending to our output array...
		clothes = append(clothes, clothing)
	})


	// start scraping at our seed address
	c.Visit(os.Getenv("SEED_ADDRESS"))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	// Dump json to the standard output
	enc.Encode(clothes)

}
