package routes

import (
	"net/http"
	"projectMod/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("web/templates/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/main", handlers.HomeHandler).Methods("GET")

	//reg
	r.HandleFunc("/registration", handlers.RegistrationPage).Methods("GET")
	r.HandleFunc("/registration", handlers.RegistrationHandler).Methods("POST")

	//movies
	r.HandleFunc("/add-movie", handlers.AddMoviePage).Methods("GET")
	r.HandleFunc("/add-movie", handlers.AddMovieHandler).Methods("POST")
	r.HandleFunc("/delete-movie", handlers.DeleteMovieHandler).Methods("POST")
	r.HandleFunc("/movie/{id}", handlers.MovieDetailHandler).Methods("GET")
	r.HandleFunc("/comment/{id}", handlers.AddCommentHandler).Methods("POST")


	//profile
	r.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	r.HandleFunc("/change-password", handlers.ChangePasswordHandler).Methods("POST")
	r.HandleFunc("/change-username", handlers.ChangeUsernameHandler).Methods("POST")


	return r
}
