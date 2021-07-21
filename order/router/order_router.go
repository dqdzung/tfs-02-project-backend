package router

import (
	"github.com/gorilla/mux"
	"project-backend/middleware"
	"project-backend/order/controller"
)

func OrderRouter(r *mux.Router) {
	r = r.PathPrefix("/orders").Subrouter()

	r.Methods("POST").Path("/").HandlerFunc(middleware.TokenAuth(controller.CreateOrder))
	r.Methods("GET").Path("/voucher/{code}").HandlerFunc(controller.GetVoucherByCode)
	r.Methods("POST").Path("/cart").HandlerFunc(middleware.TokenAuth(controller.CheckCart))

}
