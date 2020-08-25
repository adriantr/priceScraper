package Kolonial

import (
	"net/http"
	"priceScraper/Interfaces"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Kolonial struct{}

func (k Kolonial) GetDoc() []*goquery.Document {

	var docs []*goquery.Document

	response, _ := http.Get("https://kolonial.no/kategorier/")

	doc, _ := goquery.NewDocumentFromReader(response.Body)

	docs = append(docs, doc)

	return docs
}

func (k Kolonial) ReturnProducts(docs []*goquery.Document) []Interfaces.Product {
	var products []Interfaces.Product

	for _, doc := range docs {
		subs := getSubcategories(getMainCategories(doc))

		products = append(products, getProductInfo(&subs)...)
	}

	return products

}
func getMainCategories(doc *goquery.Document) []string {
	var mainCat []string

	doc.Find(".parent-category").Each(func(i int, s *goquery.Selection) {
		lol := s.Find("a")
		link, _ := lol.Attr("href")
		mainCat = append(mainCat, link)
	})

	return mainCat
}

func getSubcategories(mainCats []string) []string {
	var subCat []string
	var wg sync.WaitGroup

	cats := make(chan []string)

	for _, mainCat := range mainCats {
		wg.Add(1)
		go extractSubCats(mainCat, &wg, cats)
	}
	for i := 0; i < len(mainCats); i++ {
		res, ok := <-cats
		if ok == false {
			break
		}
		subCat = append(subCat, res...)
	}

	wg.Wait()
	close(cats)
	return subCat
}

func extractSubCats(link string, wg *sync.WaitGroup, c chan []string) {
	var subCat []string

	defer wg.Done()

	response, _ := http.Get("https://kolonial.no" + link)

	doc, _ := goquery.NewDocumentFromReader(response.Body)

	doc.Find(".child-category").Each(func(i int, s *goquery.Selection) {
		lol := s.Find("a")

		link, _ := lol.Attr("href")

		subCat = append(subCat, link)
	})

	c <- subCat
}

func getProductInfo(cats *[]string) []Interfaces.Product {
	var products []Interfaces.Product
	var wg sync.WaitGroup
	prods := make(chan []Interfaces.Product)

	for _, cat := range *cats {
		wg.Add(1)
		go visitSubCat(cat, &wg, prods)
	}
	for i := 0; i < len(*cats); i++ {
		res, ok := <-prods
		if ok == false {
			break
		}
		products = append(products, res...)
	}

	wg.Wait()
	close(prods)
	return products

}

func visitSubCat(link string, wg *sync.WaitGroup, c chan []Interfaces.Product) {
	var prods []Interfaces.Product
	defer wg.Done()

	response, _ := http.Get("https://kolonial.no" + link)

	doc, _ := goquery.NewDocumentFromReader(response.Body)

	doc.Find(".product-list-item ").Each(func(i int, s *goquery.Selection) {
		var product Interfaces.Product
		lol := s.Find(".label-price")
		product.Price = strings.TrimSpace(lol.Text())
		if len(product.Price) > 0 {
			product.Price = product.Price[3:len(product.Price)]
		} else {
			product.Price = "0"
		}

		test := s.Find(".name-main").Text()
		test1 := s.Find(".name-extra").Text()
		test2 := strings.Split(test1, "\n")

		product.Name = test + " - " + strings.TrimSpace(test2[3])

		prods = append(prods, product)
	})

	c <- prods
}
