package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"

	"github.com/gorilla/mux"
)

func MovieDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.DB.QueryRow(`SELECT id, title, description, director, release_year, genre, poster_url, trailer_url FROM movies WHERE id = $1`, id)

	var movie struct {
		ID          int
		Title       string
		Description string
		Director    string
		ReleaseYear int
		Genre       string
		PosterURL   string
		TrailerURL  string
	}

	err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Director, &movie.ReleaseYear, &movie.Genre, &movie.PosterURL, &movie.TrailerURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/movieDetail.html"))
	tmpl.Execute(w, movie)
}
