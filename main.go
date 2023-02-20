package main

import (
	"github.com/gorilla/mux"
	"log"
	"mvc/config"
	"mvc/controllers"
	"net/http"
)

func main() {
	productController := controller.NewProductController(config.NewMySQLConnection(), config.NewRedisConnection())
	router := mux.NewRouter()

	router.HandleFunc("/products", productController.GetAll).Methods("GET")
	router.HandleFunc("/products/{id}", productController.GetById).Methods("GET")
	router.HandleFunc("/products", productController.Create).Methods("POST")
	router.HandleFunc("/products/{id}", productController.Update).Methods("PUT")
	router.HandleFunc("/products/{id}", productController.Delete).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
