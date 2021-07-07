package auth

import (
	auth "project-backend/auth/controller"

	"github.com/gorilla/mux"
)

func AuthRouter(r *mux.Router) {
	r = r.PathPrefix("/auth").Subrouter()
	r.Methods("POST").Path("/login").HandlerFunc(auth.Login)
	r.Methods("POST").Path("/signup").HandlerFunc(auth.SignUp)
}
