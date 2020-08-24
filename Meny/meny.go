package Meny

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"priceScraper/Interfaces"
	"priceScraper/Util"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Meny struct{}

func (m Meny) GetDoc() []*goquery.Document {
	var docs []*goquery.Document
	i := 1

	for {

		url := "https://meny.no/WSPager?categorySlug1=&categorySlug2=&categorySlug3=&categorySlug4=&offersOnly=false&page=#&pageSize=200&query=&sort=&allergens=&fallbackBestSales=false&addHouseHoldId=false&id=&oftenBought=false&blockId=Ofte-kj%c3%b8pt+blokk"

		url = strings.Replace(url, "#", strconv.Itoa(i), 1)
		fmt.Print("Kall mot meny " + strconv.Itoa(i))
		response, _ := http.Get(url)

		bodyBytes, _ := ioutil.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		if bodyString == "" {
			break
		}
		bodyString = "<html><head></head><body>" + bodyString + "</body></html>"
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
		bodyString, _ = doc.Html()
		docs = append(docs, doc)
		i++
	}
	return docs
}

func (m Meny) ReturnProducts(docs []*goquery.Document) []Interfaces.Product {
	var links []string
	var products []Interfaces.Product
	// var wg sync.WaitGroup

	// chanProd := make(chan Interfaces.Product)

	for _, doc := range docs {
		doc.Find(".cw-product__link").Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists == true {
				if !itemExists(links, link) {
					links = append(links, link)

					res := visitProduct(link)
					if (Interfaces.Product{}) != res {
						products = append(products, res)
					} else {
						fmt.Print("Failed: " + link)
					}
				}
			}
		})
	}

	// for _, link := range links {
	// 	wg.Add(1)
	// 	go visitProduct(link, &wg, chanProd)
	// }

	// for i := 0; i < len(links); i++ {
	// 	res, ok := <-chanProd
	// 	if ok == false {
	// 		break
	// 	}
	// 	if (Interfaces.Product{}) != res {
	// 		products = append(products, res)
	// 	} else {
	// 		fmt.Print(res)
	// 	}
	// }

	return products
}

func visitProduct(link string) Interfaces.Product {

	product := Interfaces.Product{}

	response, err := http.Get(link)

	if err != nil {
		return product
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatal(err)
		fmt.Print(response.StatusCode)
	}

	//art := strings.TrimSpace(doc.Find(".cw-product-detail__title").Text())
	//txt := strings.Split(art, "\n")

	product.Name = Util.TrimArticleText(doc.Find(".cw-product-detail__title").Text())
	//string(txt[0]) + string(" - ") + string(txt[2])

	kroner := doc.Find(".cw-product__prices .cw-product__price__main").First().Text()
	ore := doc.Find(".cw-product__prices .cw-product__price__cents").First().Text()

	kronerInt, _ := strconv.Atoi(kroner)

	if kronerInt > 0 {
		product.Price = kroner + "," + ore
	} else {
		product.Price = kroner
	}

	ean := strings.Split(link, "-")

	product.Ean = ean[len(ean)-1]

	// doc.Find("meta").Each(func(i int, s *goquery.Selection) {
	// 	if name, _ := s.Attr("property"); name == "og:image" {
	// 		img, _ := s.Attr("content")
	// 		product.Img = img
	// 	}
	// })

	return product
}

func itemExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
