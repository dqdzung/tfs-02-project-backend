package product

import (
	"project-backend/product/controller"

	"github.com/gorilla/mux"
)

func ProductRouter(r *mux.Router) {
	r = r.PathPrefix("/products").Subrouter()
	r.Methods("GET").Path("/").HandlerFunc(controller.GetAll)
	// r.Methods("GET").Path("/filter=?").HandlerFunc(controller.GetAll)
	r.Methods("POST").Path("/").HandlerFunc(controller.Add)
	r.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(controller.GetOne)
	r.Methods("PUT").Path("/{id:[0-9]+}").HandlerFunc(controller.UpdateOne)
	r.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(controller.DeleteOne)
}
