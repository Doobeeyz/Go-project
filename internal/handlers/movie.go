package handlers

import(
	"html/template"
	"net/http"
	"projectMod/internal/database"
)

func AddMoviePage(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("web/templates/addMoviePage.html"))
	tmpl.Execute(w, nil)
}

func AddMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/add-movie", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	director := r.FormValue("director")
	releaseYear := r.FormValue("release_year")
	genre := r.FormValue("genre")
	posterURL := r.FormValue("poster_url")
	trailerURL := r.FormValue("trailer_url")

	_, err := database.DB.Exec(`
		INSERT INTO movies (title, description, director, release_year, genre, poster_url, trailer_url) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		title, description, director, releaseYear, genre, posterURL, trailerURL)

	if err != nil {
		http.Error(w, "Ошибка при добавлении фильма", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")

	_, err := database.DB.Exec("DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Не удалось удалить фильм", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}