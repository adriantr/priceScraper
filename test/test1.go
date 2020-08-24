package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	url := "https://meny.no/WSPager?categorySlug1=&categorySlug2=&categorySlug3=&categorySlug4=&offersOnly=false&page=#&pageSize=200&query=&sort=&allergens=&fallbackBestSales=true&addHouseHoldId=true&id=&oftenBought=true&blockId=Ofte-kj%c3%b8pt+blokk"

	url = strings.Replace(url, "#", strconv.Itoa(1), -1)

	fmt.Print(url)
	//response, _ := http.Get(url)

	//bodyBytes, _ := ioutil.ReadAll(response.Body)

	//fmt.Print(string(bodyBytes))

	//if string(bodyBytes) == "" {
	//	fmt.Print("TOM")
	//}
}

// func main() {
// 	var doc *goquery.Document
// 	response, _ := http.Get("https://kolonial.no/kategorier/")

// 	doc, _ = goquery.NewDocumentFromReader(response.Body)

// 	fmt.Print(returnProducts(doc)[0])
// }
// func returnProducts(doc *goquery.Document) []Interfaces.Product {

// 	subs := getSubcategories(getMainCategories(doc))

// 	products := getProductInfo(&subs)

// 	return products

// }
// func getMainCategories(doc *goquery.Document) []string {
// 	var mainCat []string

// 	doc.Find(".parent-category").Each(func(i int, s *goquery.Selection) {
// 		lol := s.Find("a")
// 		link, _ := lol.Attr("href")
// 		mainCat = append(mainCat, link)
// 	})

// 	return mainCat
// }

// func getSubcategories(mainCats []string) []string {
// 	var subCat []string
// 	var wg sync.WaitGroup

// 	cats := make(chan []string)

// 	for _, mainCat := range mainCats {
// 		wg.Add(1)
// 		go extractSubCats(mainCat, &wg, cats)
// 	}
// 	for i := 0; i < len(mainCats); i++ {
// 		res, ok := <-cats
// 		if ok == false {
// 			break
// 		}
// 		subCat = append(subCat, res...)
// 	}

// 	wg.Wait()
// 	close(cats)
// 	return subCat
// }

// func extractSubCats(link string, wg *sync.WaitGroup, c chan []string) {
// 	var subCat []string

// 	defer wg.Done()

// 	response, _ := http.Get("https://kolonial.no" + link)

// 	doc, _ := goquery.NewDocumentFromReader(response.Body)

// 	doc.Find(".child-category").Each(func(i int, s *goquery.Selection) {
// 		lol := s.Find("a")

// 		link, _ := lol.Attr("href")

// 		subCat = append(subCat, link)
// 	})

// 	c <- subCat
// }

// func getProductInfo(cats *[]string) []Interfaces.Product {
// 	var products []Interfaces.Product
// 	var wg sync.WaitGroup
// 	prods := make(chan []Interfaces.Product)

// 	for _, cat := range *cats {
// 		wg.Add(1)
// 		go visitSubCat(cat, &wg, prods)
// 	}
// 	for i := 0; i < len(*cats); i++ {
// 		res, ok := <-prods
// 		if ok == false {
// 			break
// 		}
// 		products = append(products, res...)
// 	}

// 	wg.Wait()
// 	close(prods)
// 	return products

// }

// func visitSubCat(link string, wg *sync.WaitGroup, c chan []Interfaces.Product) {
// 	var prods []Interfaces.Product
// 	defer wg.Done()

// 	response, _ := http.Get("https://kolonial.no" + link)

// 	doc, _ := goquery.NewDocumentFromReader(response.Body)

// 	doc.Find(".product-list-item ").Each(func(i int, s *goquery.Selection) {
// 		var product Interfaces.Product
// 		lol := s.Find(".label-price")
// 		product.Price = lol.Text()

// 		product.Name = strings.TrimSpace(s.Find(".name-main").Text()) + " - " + strings.TrimSpace(s.Find(".name-extra").Text())

// 		prods = append(prods, product)
// 	})

// 	c <- prods
// }
