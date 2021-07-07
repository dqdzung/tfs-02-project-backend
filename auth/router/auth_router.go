package auth

import (
	"project-backend/auth/controller"

	"github.com/gorilla/mux"
)

func AuthRouter(r *mux.Router) {
	r = r.PathPrefix("/auth").Subrouter()
	r.Methods("POST").Path("/login").HandlerFunc(controller.Login)
	r.Methods("POST").Path("/signup").HandlerFunc(controller.SignUp)
}
