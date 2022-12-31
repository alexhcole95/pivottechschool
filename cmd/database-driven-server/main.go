package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	databasePath = "./seeder/products.db"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func pingDB(w http.ResponseWriter) bool {
	if err := db.Ping(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	return false
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}

	query := "select * from products"
	params := r.URL.Query()
	sort := params.Get("column")
	if sort != "" {
		switch sort {
		case "id", "name", "price":
			query += fmt.Sprintf(" order by %s asc", sort)
		default:
			log.Printf("error: %s invalid argument for column\n", sort)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	limit := params.Get("limit")
	if limit != "" {
		limInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Printf("error: %s invalid argument for limit\n", limit)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query += fmt.Sprintf(" limit %d", limInt)
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var prods []Product
	for rows.Next() {
		var prod Product
		if err = rows.Scan(&prod.ID, &prod.Name, &prod.Price); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		prods = append(prods, prod)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(prods)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	if _, err := w.Write(bs); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("select * from products where id=%d", idInt)
	var prod Product
	if err := db.QueryRow(query).Scan(&prod.ID, &prod.Name, &prod.Price); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	bs, err := json.Marshal(prod)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	if _, err := w.Write(bs); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func createNewProductHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var prod Product
	if err := json.Unmarshal(body, &prod); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if prod.ID != 0 || len(strings.Fields(prod.Name)) == 0 || prod.Price == 0 {
		log.Println("error: invalid json payload from post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := "insert into products(name, price) values(?, ?)"
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := stmt.Exec(prod.Name, prod.Price); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}

	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var oldProd Product
	if err := db.QueryRow("select * from products where id=?", idInt).Scan(&oldProd.ID, &oldProd.Name, &oldProd.Price); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var newProd Product
	if err := json.Unmarshal(body, &newProd); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (idInt != newProd.ID && newProd.ID != 0) || len(strings.Fields(newProd.Name)) == 0 || newProd.Price == 0 {
		log.Println("error: invalid json payload from put request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := "update products set name = ?, price = ? where id = ?"
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := stmt.Exec(newProd.Name, newProd.Price, idInt); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}

	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var prod Product
	if err := db.QueryRow("select * from products where id=?", idInt).Scan(&prod.ID, &prod.Name, &prod.Price); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	query := "delete from products where id=?"
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := stmt.Exec(idInt); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	database, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	db = database

	router := mux.NewRouter()
	router.HandleFunc("/products", getProductsHandler).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", getProductByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/products", createNewProductHandler).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", updateProductHandler).Methods(http.MethodPut)
	router.HandleFunc("/products/{id}", deleteProductHandler).Methods(http.MethodDelete)

	log.Println("Starting server. Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
