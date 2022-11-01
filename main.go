package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var products []product

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func returnProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "404", http.StatusInternalServerError)
	}

	for _, product := range products {
		if int64(product.ID) == id {
			if err = json.NewEncoder(w).Encode(product); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func createNewProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var item product
	if err := json.Unmarshal(reqBody, &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	products = append(products, item)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "404", http.StatusInternalServerError)
	}

	var updatedProduct product
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err = json.Unmarshal(reqBody, &updatedProduct); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	for i, product := range products {
		if int64(product.ID) == id {
			product.ID = updatedProduct.ID
			product.Name = updatedProduct.Name
			product.Description = updatedProduct.Description
			product.Price = updatedProduct.Price
			products[i] = product
			if err = json.NewEncoder(w).Encode(product); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "404", http.StatusInternalServerError)
	}

	for index, product := range products {
		if int64(product.ID) == id {
			products = append(products[:index], products[index+1:]...)
		}
	}
}

func initProducts() {
	bs, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(bs, &products); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initProducts()
	r := mux.NewRouter()

	r.HandleFunc("/products", getProducts)
	r.HandleFunc("/product", createNewProduct).Methods("POST")
	r.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/product/{id}", returnProduct)

	log.Println("Listening on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
