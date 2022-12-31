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

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}

	for _, product := range products {
		if int64(product.ID) == id {
			if err = json.NewEncoder(w).Encode(product); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func createNewProductHandler(w http.ResponseWriter, r *http.Request) {
	var item product
	var items []product
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", 400)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(reqBody, &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	items = append(products, item)
	bs, err := json.Marshal(items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := os.WriteFile("products.json", bs, 0666); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	products = items
	w.WriteHeader(http.StatusCreated)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}

	var updatedProduct product
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", 400)
	}
	defer r.Body.Close()

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
			if err = json.NewEncoder(w).Encode(http.StatusOK); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	var deletedItem bool
	var newProduct []product
	vars := mux.Vars(r)

	for _, v := range products {
		if strconv.Itoa(v.ID) != vars["id"] {
			newProduct = append(newProduct, v)
		} else {
			deletedItem = true
		}
	}

	if deletedItem {
		bs, err := json.Marshal(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if err := os.WriteFile("products.json", bs, 0666); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	products = newProduct
	w.WriteHeader(http.StatusOK)
}

func initProducts() {
	byteSlice, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(byteSlice, &products); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initProducts()
	r := mux.NewRouter()

	r.HandleFunc("/products", getProductsHandler).Methods("GET")
	r.HandleFunc("/products/{id}", getProductByIDHandler).Methods("GET")
	r.HandleFunc("/products", createNewProductHandler).Methods("POST")
	r.HandleFunc("/products/{id}", updateProductHandler).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProductHandler).Methods("DELETE")

	log.Println("Listening on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
