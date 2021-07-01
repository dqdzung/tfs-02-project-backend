package server

import (
	"fmt"
	"net/http"
	"project-backend/controller"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func RunServer() {
	const PORT = 8000
	fmt.Printf("Server running at port %v...", PORT)
	defer fmt.Println("Server stopped!")
	router := mux.NewRouter().StrictSlash(true)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	handler := c.Handler(router)

	router.Methods("GET").Path("/products").HandlerFunc(controller.GetAllProducts) //done
	router.Methods("POST").Path("/products").HandlerFunc(controller.AddProduct)    //done

	router.Methods("GET").Path("/products/{id:[0-9]+}").HandlerFunc(controller.GetOneProduct)
	router.Methods("PUT").Path("/products/{id:[0-9]+}").HandlerFunc(controller.UpdateProduct)
	router.Methods("DELETE").Path("/products/{id:[0-9]+}").HandlerFunc(controller.DeleteProduct)

	router.Methods("GET").Path("/orders").HandlerFunc(controller.GetAllOrders)
	router.Methods("POST").Path("/orders").HandlerFunc(controller.AddOrder)

	router.Methods("GET").Path("/orders/{id:[0-9]+}").HandlerFunc(controller.GetOrder)
	router.Methods("PUT").Path("/orders/{id:[0-9]+}").HandlerFunc(controller.UpdateOrder)
	router.Methods("DELETE").Path("/orders/{id:[0-9]+}").HandlerFunc(controller.DeleteOrder)

	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		panic(err)
	}
}
