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

func RegisterMeasureRoutes(r *mux.Router, db *pgxpool.Pool) {
	r.HandleFunc("/measure/", getMeasures(db)).Methods("GET")
	r.HandleFunc("/measure/{id}", getMeasureByID(db)).Methods("GET")
	r.HandleFunc("/measure/", createMeasure(db)).Methods("POST")
	r.HandleFunc("/measure/{id}", updateMeasure(db)).Methods("PUT")
	r.HandleFunc("/measure/{id}", deleteMeasure(db)).Methods("DELETE")
}

func getMeasures(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(context.Background(), "SELECT id, name FROM measures")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var measures []models.Measure
		for rows.Next() {
			var m models.Measure
			if err := rows.Scan(&m.ID, &m.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			measures = append(measures, m)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(measures)
	}
}

func getMeasureByID(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid measure ID", http.StatusBadRequest)
			return
		}

		var m models.Measure
		err = db.QueryRow(context.Background(), "SELECT id, name FROM measures WHERE id=$1", id).Scan(&m.ID, &m.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Measure not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	}
}

func createMeasure(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m models.Measure
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var id int
		err := db.QueryRow(context.Background(), "INSERT INTO measures (name) VALUES ($1) RETURNING id", m.Name).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}
}

func updateMeasure(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid measure ID", http.StatusBadRequest)
			return
		}

		var m models.Measure
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "UPDATE measures SET name=$1 WHERE id=$2", m.Name, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteMeasure(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid measure ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "DELETE FROM measures WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
