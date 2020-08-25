package main

import (
	"priceScraper/Coop"
	"priceScraper/Interfaces"
	"priceScraper/Kolonial"
	"priceScraper/Meny"

	"github.com/sahilm/fuzzy"
)

func main() {
	allProducts := Interfaces.AllProducts{}
	meny := Meny.Meny{}
	coop := Coop.Coop{}
	kolonial := Kolonial.Kolonial{}

	menyDoc := Interfaces.HTMLScraper.GetDoc(meny)
	allProducts.Meny = Interfaces.HTMLScraper.ReturnProducts(meny, menyDoc)

	coopCat := Interfaces.APIScraper.GetCategories(coop)
	allProducts.Coop = Interfaces.APIScraper.GetProducts(coop, coopCat)

	kolonialDoc := Interfaces.HTMLScraper.GetDoc(kolonial)
	allProducts.Kolonial = Interfaces.HTMLScraper.ReturnProducts(kolonial, kolonialDoc)

	guessKolonialGtin(&allProducts)

	//AZBus.SendToAzure(&allProducts)
}

type products []Interfaces.Product

func (p products) String(i int) string {
	return p[i].Name
}
func (p products) Len() int {
	return len(p)
}

func guessKolonialGtin(prods *Interfaces.AllProducts) {
	var all products

	all = append(all, prods.Coop...)
	all = append(all, prods.Meny...)

	for _, product := range prods.Kolonial {
		results := fuzzy.FindFrom(product.Name, all)
		if len(results) > 0 {
			product.Ean = all[results[0].Index].Ean
		}
	}

}
