package routes

import (
	"projectMod/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/main", handlers.HomeHandler).Methods("GET")

	//reg
	r.HandleFunc("/registration", handlers.RegistrationPage).Methods("GET")
	r.HandleFunc("/registration", handlers.RegistrationHandler).Methods("POST")
	return r
}
