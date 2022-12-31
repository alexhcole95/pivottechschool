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

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

const (
	productsPath = "./seeder/products.db"
)

var db *sql.DB

func main() {
	database, err := sql.Open("sqlite3", productsPath)
	if err != nil {
		log.Fatal(err)
	}
	db = database
	router := mux.NewRouter()
	router.HandleFunc("/products", getProdsHandler).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", getProdByIDHandler).Methods(http.MethodGet)
	router.HandleFunc("/products", addProdHandler).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", updateProdByIDHandler).Methods(http.MethodPut)
	router.HandleFunc("/products/{id}", deleteProdByIDHandler).Methods(http.MethodDelete)
	log.Println("Starting server. Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getProdsHandler(w http.ResponseWriter, r *http.Request) {
	if pingDB(w) {
		return
	}
	query := "select * from products"
	params := r.URL.Query()
	sort := params.Get("column")
	if sort != "" {
		switch sort {
		case "id", "name", "description", "price":
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
		if err = rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price); err != nil {
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

func getProdByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := db.QueryRow(query).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price); err != nil {
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

func addProdHandler(w http.ResponseWriter, r *http.Request) {
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
	if prod.ID != 0 || len(strings.Fields(prod.Name)) == 0 || len(strings.Fields(prod.Description)) == 0 || prod.Price == 0 {
		log.Println("error: invalid json payload from post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := "insert into products(name, description, price) values(?, ?, ?)"
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
	if _, err := stmt.Exec(prod.Name, prod.Description, prod.Price); err != nil {
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

func updateProdByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := db.QueryRow("select * from products where id=?", idInt).Scan(&oldProd.ID, &oldProd.Name, &oldProd.Description, &oldProd.Price); err != nil {
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
	if (idInt != newProd.ID && newProd.ID != 0) || len(strings.Fields(newProd.Name)) == 0 || len(strings.Fields(newProd.Description)) == 0 || newProd.Price == 0 {
		log.Println("error: invalid json payload from put request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := "update products set name = ?, description = ?, price = ? where id = ?"
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
	if _, err := stmt.Exec(newProd.Name, newProd.Description, newProd.Price, idInt); err != nil {
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

func deleteProdByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := db.QueryRow("select * from products where id=?", idInt).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price); err != nil {
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

func pingDB(w http.ResponseWriter) bool {
	if err := db.Ping(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	return false
}
