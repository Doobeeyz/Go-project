package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/models"

	"github.com/gorilla/mux"
)

func AddToFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)
	movieID := mux.Vars(r)["id"]

	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO favorites (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, movieID)
	if err != nil {
		http.Error(w, "Ошибка добавления в любимое", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/movie/"+movieID, http.StatusSeeOther)
}

func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)

	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(`
		SELECT m.id, m.title, m.release_year, m.genre, m.description, m.poster_url
		FROM favorites f
		JOIN movies m ON f.movie_id = m.id
		WHERE f.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, "Ошибка получения избранных фильмов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favorites []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Genre, &m.Description, &m.PosterURL)
		if err != nil {
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
			return
		}
		favorites = append(favorites, m)
	}

	roleID, _ := session.Values["role_id"].(int)

	tmpl, err := template.ParseFiles("web/templates/favorites.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		Favorites []models.Movie
		RoleID    int
	}{
		Favorites: favorites,
		RoleID:    roleID,
	}

	tmpl.Execute(w, data)
}

func RemoveFromFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)
	movieID := mux.Vars(r)["id"]

	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("DELETE FROM favorites WHERE user_id = $1 AND movie_id = $2", userID, movieID)
	if err != nil {
		http.Error(w, "Ошибка удаления из любимого", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}
