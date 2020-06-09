package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

type Product struct {
	Uuid string `json:"uuid"`
	Product string `json:"product"`
	Price float64 `json:"price, string"`
}

type Products struct {
	Product []Product
}

func loadData() []byte {
	jsonFile, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	return data
}

func ListProducts(w http.ResponseWriter, r * http.Request) {
	products := loadData();
	w.Write([]byte(products))
}

func getProductById(w http.ResponseWriter, r * http.Request) {
	vars := mux.Vars(r)
	data := loadData()

	products := loadData();
	w.Write([]byte(products))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", ListProducts)
	http.ListenAndServe(":8081", r)
}
