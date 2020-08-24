package Interfaces

import (
	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	Ean   string
	Name  string
	Price string
	Img   string
}

type AllProducts struct {
	Meny     []Product
	Coop     []Product
	Kolonial []Product
}

type HTMLScraper interface {
	GetDoc() []*goquery.Document
	ReturnProducts(docs []*goquery.Document) []Product
}

type APIScraper interface {
	GetCategories() *[]string
	GetProducts(categories *[]string) []Product
}
