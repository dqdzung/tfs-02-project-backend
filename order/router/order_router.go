package router

import (
	"project-backend/middleware"
	"project-backend/order/controller"

	"github.com/gorilla/mux"
)

func OrderRouter(r *mux.Router) {
	r = r.PathPrefix("/orders").Subrouter()

	r.Methods("POST").Path("").HandlerFunc(middleware.TokenAuth(controller.CreateOrder))
	r.Methods("GET").Path("/voucher/{code}").HandlerFunc(controller.GetVoucherByCode)
	r.Methods("POST").Path("/cart").HandlerFunc(middleware.TokenAuth(controller.CheckCart))

}
