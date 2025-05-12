package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/models"
)

type PageData struct {
	RoleID int
	Movies []models.Movie
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)

	var roleID int
	_ = database.DB.QueryRow("SELECT role_id FROM users WHERE username = $1", username).Scan(&roleID)

	rows, err := database.DB.Query("SELECT id, title, description, director, release_year, genre, poster_url FROM movies")
	if err != nil {
		http.Error(w, "Ошибка получения фильмов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Director, &m.ReleaseYear, &m.Genre, &m.PosterURL); err == nil {
			movies = append(movies, m)
		}
	}

	data := PageData{
		RoleID: roleID,
		Movies: movies,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/mainPage.html"))
	tmpl.Execute(w, data)
}
