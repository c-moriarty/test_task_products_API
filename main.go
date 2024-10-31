package main

import (
	"context"
	"log"
	"net/http"
	"product_app/handlers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:1@localhost:5432/products")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	handlers.RegisterProductRoutes(r, db)
	handlers.RegisterMeasureRoutes(r, db)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
