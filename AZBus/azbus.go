package AZBus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"priceScraper/Interfaces"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type Split struct {
	Chain    string
	Products []Interfaces.Product
}

func SendToAzure(products *Interfaces.AllProducts) {
	cstring := os.Getenv("CONNECTIONSTRING")

	cstringbyte, _ := base64.StdEncoding.DecodeString(cstring)

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(string(cstringbyte)))
	if err != nil {
		return
	}

	client, err := ns.NewQueue("prices")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := servicebus.Message{}

	message.ContentType = "application/json"

	chunks := splitIntoChunks(products)

	for _, chunk := range chunks {

		jsonData, err := json.Marshal(chunk)
		if err != nil {
			fmt.Println(err)
			return
		}
		message.Data = jsonData

		if err := client.Send(ctx, &message); err != nil {
			fmt.Println("FATAL: ", err)
		}
	}

}

func splitIntoChunks(products *Interfaces.AllProducts) []Split {
	var wg sync.WaitGroup
	var tot []Split

	prds := make(chan []Split, 3)

	wg.Add(3)

	go splitProducts(&products.Coop, "Coop", &wg, prds)
	go splitProducts(&products.Meny, "Meny", &wg, prds)
	go splitProducts(&products.Kolonial, "Kolonial", &wg, prds)

	for i := 0; i < 3; i++ {
		res, ok := <-prds
		if ok == false {
			break
		}
		tot = append(tot, res...)

	}

	wg.Wait()
	close(prds)

	return tot

}

func splitProducts(products *[]Interfaces.Product, chain string, wg *sync.WaitGroup, c chan []Split) {
	defer wg.Done()
	var splittot []Split

	for i := 0; i < len(*products); i += 1000 {
		var split Split

		split.Chain = chain
		split.Products = append(split.Products, (*products)[i:min(i+1000, len(*products))]...)
		splittot = append(splittot, split)
	}

	c <- splittot
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
