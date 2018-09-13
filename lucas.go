package main

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"strings"
	"github.com/gocolly/colly"
	"github.com/fatih/color"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
)

type Clothing struct {
	Name					string
	Code					string
	Description		string
	Price					float64
}

func dbWrite(product Clothing) {
	const (
	  host     = "localhost"
	  port     = 5432
	  user     = "user"
	  // password = ""
	  dbname   = "lucas_db"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable",
    host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  log.Print("Successfully connected!")
	fmt.Printf("%s, %s, %s, %f", product.Name, product.Code, product.Description, product.Price)
	sqlStatement := `
	INSERT INTO floryday (product, code, description, price)
	VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, product.Name, product.Code, product.Description, product.Price)
	if err != nil {
	  panic(err)
	}
}

func main() {

	c := colly.NewCollector(
		// colly.AllowedDomains("https://www.floryday.com/"),
		colly.CacheDir(".floryday_cache"),
  	// colly.MaxDepth(5), // keeping crawling limited for our initial experiments
  )

	// clothing detail scraping collector
	detailCollector := c.Clone()

	clothes := make([]Clothing, 0, 200)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

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

		// links provided need to be better filtered
		// hardcoding one value only to work here for now...
		if strings.Contains(clothingURL, "-Dress-"){
			// Activate detailCollector
			color.Green("Crawling Link Validated -> Commencing Crawl for %s", clothingURL)
			detailCollector.Visit(clothingURL)
		} else {
			color.Red("Validation Failed -> Cancelling Crawl for %s", clothingURL)
			return
		}

	})

	// Extract details of the clothing
	detailCollector.OnHTML(`div[class=prod-right-in]`, func(e *colly.HTMLElement) {
		// TODO secure variables with default error strings in event values are missing
		title := e.ChildText(".prod-name")
		code := strings.Split(e.ChildText(".prod-item-code"), "#")[1]
		stringPrice := strings.TrimPrefix(e.ChildText(".prod-price"),"€ ") // TODO non scalable outside eurozone
		price, err := strconv.ParseFloat(stringPrice, 64) // conversion to float64
		color.Red("err in parsing price -> %s", err)
		description := e.ChildText(".grid-uniform")

		clothing := Clothing{
			Name: 					title,
			Code: 					code,
			Description: 		description,
			Price:					price,
		}

		// writing as we go to DB
		// TODO optiize to handle bulk array uplaods instead
		dbWrite(clothing)

		// appending to our output array...
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
