package Coop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"priceScraper/Interfaces"

	"github.com/tidwall/gjson"
)

type Coop struct{}

func (c Coop) GetCategories() *[]string {
	var categories []string

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://matlevering.coop.no/", nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("x-includeappshelldata", "true")

	resp, _ := client.Do(req)

	bd, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bd)

	if !gjson.Valid(bodyString) {
		fmt.Print("Invalid json")
	}

	menu := gjson.Get(bodyString, "appShellData.mainMenu.mainMenuItems")

	menu.ForEach(func(key, value gjson.Result) bool {
		//println(value.String())
		children := value.Get("children")

		if children.Raw == "[]" {
			categories = append(categories, value.Get("url").String())
		} else {

			children.ForEach(func(key, child gjson.Result) bool {
				categories = append(categories, child.Get("url").String())
				return true
			})
		}
		return true
	})

	return &categories
}

func (c Coop) GetProducts(categories *[]string) []Interfaces.Product {
	var products []Interfaces.Product
	for _, cat := range *categories {
		client := &http.Client{}

		req, _ := http.NewRequest("GET", "https://matlevering.coop.no"+cat, nil)
		req.Header.Add("x-requested-with", "XMLHttpRequest")
		req.Header.Add("x-includeappshelldata", "true")

		resp, _ := client.Do(req)

		bd, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bd)

		prds := gjson.Get(bodyString, "products")

		prds.ForEach(func(key, prod gjson.Result) bool {
			var product Interfaces.Product
			product.Price = prod.Get("price.current.inclVat").String()
			product.Ean = prod.Get("variationCode").String()
			product.Img = prod.Get("images.0.url").String()
			product.Name = prod.Get("displayName").String()

			products = append(products, product)
			return true
		})
	}

	return products

}
