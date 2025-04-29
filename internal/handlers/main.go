package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
)

type Movie struct {
	ID          int
	Title       string
	Description string
	Director    string
	ReleaseYear int
	Genre       string
	PosterURL   string
	TrailerURL  string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")

	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Доступ запрещен. Войдите в систему.", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`SELECT id, title, description, director, release_year, genre, poster_url, trailer_url FROM movies`)
	
	if err != nil {
		http.Error(w, "Ошибка при загрузке фильмов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next(){
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Director, &m.ReleaseYear, &m.Genre, &m.PosterURL, &m.TrailerURL)
		if err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		movies = append(movies, m)
	}
	tmpl := template.Must(template.ParseFiles("web/templates/mainPage.html"))
	tmpl.Execute(w, movies)
}


