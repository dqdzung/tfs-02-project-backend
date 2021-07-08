package product

import (
	product "project-backend/product/controller"

	"github.com/gorilla/mux"
)

func ProductRouter(r *mux.Router) {
	r = r.PathPrefix("/products").Subrouter()
	r.Methods("GET").Path("/").HandlerFunc(product.GetAll)
	// r.Methods("GET").Path("/filter=?").HandlerFunc(controller.GetAll)
	r.Methods("POST").Path("/").HandlerFunc(product.Add)
	r.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(product.GetOne)
	r.Methods("PUT").Path("/{id:[0-9]+}").HandlerFunc(product.UpdateOne)
	r.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(product.DeleteOne)
}
