package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"reflect"
	"strconv"
	"productScraper"
)

type Product struct {
	ean string
	name string
	price string
	img string
}

var links []string;
var products []Product;

func main() {

	response, err := http.Get("https://meny.no/WSPager?categorySlug1=&categorySlug2=&categorySlug3=&categorySlug4=&offersOnly=false&page=1&pageSize=10&query=&sort=&allergens=&fallbackBestSales=true&addHouseHoldId=true&id=&oftenBought=true&blockId=Ofte-kj%c3%b8pt+blokk")

	if err != nil {
		fmt.Print("Naaah")
	} else {

		doc, err := goquery.NewDocumentFromReader(response.Body);

		if err != nil {
			log.Fatal(err);
		}

		doc.Find(".cw-product__link").Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists == true {

				if !itemExists(links, link) {
					links = append(links, link)
					visitProduct(link)
				}
			} else {
				fmt.Print("no link")
			}
			
		})

		fmt.Print(products)

	}
}

func visitProduct (link string) {
	var product Product

	response, err := http.Get(link)

	if err != nil {
		fmt.Print("Cannot access product")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body);

	if err != nil {
		log.Fatal(err);
	}

	art := strings.TrimSpace(doc.Find(".cw-product-detail__title").Text())
	txt := strings.Split(art, "\n")
	product.name = string(txt[0]) + string(" - ") + string(txt[2])

	kroner := doc.Find(".cw-product__prices .cw-product__price__main").First().Text()
	ore := doc.Find(".cw-product__prices .cw-product__price__cents").First().Text();

	kronerInt, _ := strconv.Atoi(kroner)

	if kronerInt > 0 {
		product.price = kroner+","+ore
	} else {
		product.price = kroner
	}

	ean := strings.Split(link, "-");

	product.ean = ean[len(ean)-1]

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("property"); name == "og:image" {
			img, _ := s.Attr("content")
			product.img = img
		}
	})

	products = append(products, product)

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