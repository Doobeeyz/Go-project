package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/models"
)

type PageData struct {
	RoleID int
	Movies []models.Movie
	Search string
	Genre  string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)

	var roleID int
	_ = database.DB.QueryRow("SELECT role_id FROM users WHERE username = $1", username).Scan(&roleID)

	// Чтение параметров из формы поиска
	search := r.URL.Query().Get("search")
	genre := r.URL.Query().Get("genre")

	query := "SELECT id, title, description, director, release_year, genre, poster_url FROM movies WHERE 1=1"
	var args []interface{}

	//фильтр по ключевым словам
	if search != "" {
		query += " AND (LOWER(title) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1))"
		args = append(args, "%"+search+"%")
	}

	//фильтр по жанру
	if genre != "" {
		paramIndex := len(args) + 1
		query += fmt.Sprintf(" AND LOWER(genre) LIKE LOWER($%d)", paramIndex)
		args = append(args, "%"+genre+"%")
	}

	//выполнение запроса
	rows, err := database.DB.Query(query, args...)
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
		Search: search,
		Genre:  genre,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/mainPage.html"))
	tmpl.Execute(w, data)
}
