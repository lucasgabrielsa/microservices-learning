package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"order.go/db"
	"order.go/queue"
	"fmt"
	"os"
	"time"
	"github.com/nu7hatch/gouuid"
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

func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4();
	order.Uuid = uuid.String()
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)
	return order
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()
	error := connection.Set(order.Uuid, string(json), 0).Err()
	if error != nil {
		panic(error.Error())
	}

	val, err := connection.Get(order.Uuid).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(order.Uuid, val)
}

func getProductById(id string) Product {
	response, error := http.Get(productsUrl + "/product/"+ id)
	if error != nil {
		fmt.Printf("The HTTP request failed with error %s/n", error)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	json.Unmarshal(data, &product)
	return product
}

func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.StartConsuming("checkout_queue", connection, in)

	for payload := range in {
		notifyOrderCreated(createOrder(payload), connection)
		fmt.Println(string(payload))
	}
}

func notifyOrderCreated(order Order, channel *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, "order_ex", "", channel);
}
