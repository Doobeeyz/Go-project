package routes

import (
	"net/http"
	"projectMod/internal/handlers"
	"projectMod/internal/middleware"

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
	//only for admin and super-admin
	r.Handle("/add-movie", middleware.RequireMinRoleID(2, http.HandlerFunc(handlers.AddMoviePage))).Methods("GET")
	r.Handle("/add-movie", middleware.RequireMinRoleID(2, http.HandlerFunc(handlers.AddMovieHandler))).Methods("POST")

	r.HandleFunc("/delete-movie", handlers.DeleteMovieHandler).Methods("POST")
	r.HandleFunc("/movie/{id}", handlers.MovieDetailHandler).Methods("GET")
	r.HandleFunc("/comment/{id}", handlers.AddCommentHandler).Methods("POST")

	//favorites
	r.HandleFunc("/favorites", handlers.FavoritesHandler).Methods("GET")
	r.HandleFunc("/favorites/add/{id}", handlers.AddToFavoritesHandler).Methods("POST")
	r.HandleFunc("/favorites/remove/{id}", handlers.RemoveFromFavoritesHandler).Methods("POST")

	//profile
	r.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	r.HandleFunc("/change-password", handlers.ChangePasswordHandler).Methods("POST")
	r.HandleFunc("/change-username", handlers.ChangeUsernameHandler).Methods("POST")

	r.Handle("/admin/roles", middleware.RequireRoleID(3, http.HandlerFunc(handlers.SuperAdminRoleHandler))).Methods("GET")
	r.HandleFunc("/admin/roles/change", handlers.ChangeUserRoleHandler).Methods("POST")

	return r
}
