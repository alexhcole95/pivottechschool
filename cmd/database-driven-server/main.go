package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var info []product

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var db *sql.DB

// item.Name == "" || item.ID == 0 || item.Price == 0
func (p product) IsValid() bool {
	if p.Name == "" {
		return false
	}
	if p.ID == 0 {
		return false
	}
	if p.Price == 0 {
		return false
	}

	return true
}

func initProducts() *sql.DB {
	var err error
	log.Println("Connecting to DB...")

	db, err = sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Println("Could not connect to DB!")
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to reach database: %v", err)
	}
	log.Println("Connected to DB!")

	return db
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving request from client...")

	lim := r.FormValue("limit")
	l, err := strconv.Atoi(lim)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}
	sort := r.FormValue("sort")

	log.Println("Request received! Searching for products...")
	rows, err := db.Query("SELECT * FROM products ORDER BY ? LIMIT ?", sort, l)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var id int
		var name string
		var price float64
		err = rows.Scan(&id, &name, &price)
		if err != nil {
			log.Fatal(err)
		}

		info = append(info, product{ID: id, Name: name, Price: price})
	}

	log.Println("Products found! Returning products to client...")
	if err = json.NewEncoder(w).Encode(info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving request from client...")

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		var price float64
		err = rows.Scan(&id, &name, &price)
		if err != nil {
			log.Fatal(err)
		}

		info = append(info, product{ID: id, Name: name, Price: price})
	}

	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}

	log.Println("Request received! Searching for product...")
	if itemID < 0 || itemID > int64(len(info)) {
		log.Println("Product not found! Returning status to client...")
		http.Error(w, "ID not found", 404)
	} else {
		for _, product := range info {
			if int64(product.ID) == itemID {
				log.Println("Product found! Returning product to client...")
				if err = json.NewEncoder(w).Encode(product); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving request from client...")

	var item product
	var items []product
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", 400)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(reqBody, &item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	items = append(info, item)
	info = items

	log.Println("Request received! Adding product to DB...")
	if item.IsValid() == false {
		http.Error(w, "Missing required fields", 400)
	} else {
		if _, err = db.Exec("INSERT INTO products (id, name, price) VALUES (?, ?, ?)", item.ID, item.Name, item.Price); err != nil {
			log.Println(err)
		}

		log.Println("Product added to DB! Returning status to client...")
		w.WriteHeader(http.StatusCreated)
	}
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving request from client...")

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var id int
		var name string
		var price float64
		err = rows.Scan(&id, &name, &price)
		if err != nil {
			log.Println(err)
		}

		info = append(info, product{ID: id, Name: name, Price: price})
	}

	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}

	log.Println("Request received! Searching for product...")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", 400)
	}
	defer r.Body.Close()

	var updatedProduct product
	if err = json.Unmarshal(reqBody, &updatedProduct); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if updatedProduct.Name == "" || updatedProduct.ID == 0 || updatedProduct.Price == 0 {
		http.Error(w, "Missing required fields", 400)
	} else {
		log.Println("Product found! Updating product in DB...")
		for i, product := range info {
			if int64(product.ID) == itemID {
				product.ID = updatedProduct.ID
				product.Name = updatedProduct.Name
				product.Price = updatedProduct.Price
				info[i] = product
				if _, err = db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", updatedProduct.Name, updatedProduct.Price, updatedProduct.ID); err != nil {
					log.Println(err)
				}
			}
		}
	}

	log.Println("Product updated in DB! Returning status to client...")
	w.WriteHeader(http.StatusCreated)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving request from client...")
	var newProduct []product

	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Unable to convert string to int", 400)
	}

	if itemID < 0 || itemID > int64(len(info)) {
		log.Println("Product not found! Returning status to client...")
		http.Error(w, "ID not found", 404)
	} else {
		log.Println("Request received! Searching for product...")
		for _, product := range info {
			if int64(product.ID) != itemID {
				newProduct = append(newProduct, product)
			}
		}
		info = newProduct

		log.Println("Product found! Deleting product from DB...")
		if _, err = db.Exec("DELETE FROM products WHERE id = ?", itemID); err != nil {
			log.Println(err)
		}

		log.Println("Product deleted from DB! Returning status to client...")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	log.Println("Starting server...")
	initProducts()

	r := mux.NewRouter()
	r.HandleFunc("/products", getProductsHandler).Methods("GET")
	r.HandleFunc("/products/{id}", getProductHandler).Methods("GET")
	r.HandleFunc("/products", addProductHandler).Methods("POST")
	r.HandleFunc("/products/{id}", updateProductHandler).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProductHandler).Methods("DELETE")

	log.Println("Listening on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
