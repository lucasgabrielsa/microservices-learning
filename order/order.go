package main

import (
	"order.go/queue"
	"fmt"
	"os"
	"time"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Order struct {
	Uuid string  `json:"uuid"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	ProductId string `json:"product_id"`
	Status string  `json:"status"`
	CreatedAt time.Time  `json:"created_at, string"`
}


var productsUrl string
func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
	fmt.Println("ProductsUrl="+ productsUrl)
}

func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.StartConsuming(connection, in)

	for payload := range in {
		fmt.Println(string(payload))
	}
}
