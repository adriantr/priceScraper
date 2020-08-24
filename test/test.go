package main

/*
// func main() {
// 	client := &http.Client{}

// 	req, _ := http.NewRequest("GET", "https://matlevering.coop.no/", nil)
// 	req.Header.Add("x-requested-with", "XMLHttpRequest")
// 	req.Header.Add("x-includeappshelldata", "true")

// 	resp, _ := client.Do(req)

// 	bd, _ := ioutil.ReadAll(resp.Body)
// 	bodyString := string(bd)

// 	if !gjson.Valid(bodyString) {
// 		fmt.Print("Invalid json")
// 	}

// 	menu := gjson.Get(bodyString, "appShellData.mainMenu.mainMenuItems")

// 	menu.ForEach(func(key, value gjson.Result) bool {
// 		//println(value.String())
// 		children := value.Get("children")

// 		if children.Raw == "[]" {
// 			println(value.Get("url").String())
// 		} else {

// 			children.ForEach(func(key, child gjson.Result) bool {
// 				println(child.Get("url").String())
// 				return true
// 			})
// 		}
// 		return true
// 	})
// }
func main() {
	var products []Interfaces.Product

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://matlevering.coop.no/hus-og-hjem/batterier", nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

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
*/
