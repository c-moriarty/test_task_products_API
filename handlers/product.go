package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"product_app/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RegisterProductRoutes(r *mux.Router, db *pgxpool.Pool) {
	r.HandleFunc("/product/", getProducts(db)).Methods("GET")
	r.HandleFunc("/product/{id}", getProductByID(db)).Methods("GET")
	r.HandleFunc("/product/", createProduct(db)).Methods("POST")
	r.HandleFunc("/product/{id}", updateProduct(db)).Methods("PUT")
	r.HandleFunc("/product/{id}", deleteProduct(db)).Methods("DELETE")
}

func getProducts(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(context.Background(), "SELECT id, name, quantity, unit_cost, measure FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.UnitCost, &p.Measure); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func getProductByID(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var p models.Product
		err = db.QueryRow(context.Background(), "SELECT id, name, quantity, unit_cost, measure FROM products WHERE id=$1", id).Scan(&p.ID, &p.Name, &p.Quantity, &p.UnitCost, &p.Measure)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func createProduct(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var id int
		err := db.QueryRow(context.Background(), "INSERT INTO products (name, quantity, unit_cost, measure) VALUES ($1, $2, $3, $4) RETURNING id", p.Name, p.Quantity, p.UnitCost, p.Measure).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}
}

func updateProduct(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "UPDATE products SET name=$1, quantity=$2, unit_cost=$3, measure=$4 WHERE id=$5", p.Name, p.Quantity, p.UnitCost, p.Measure, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteProduct(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "DELETE FROM products WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
