package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/models"

	"github.com/gorilla/mux"
)

func MovieDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.DB.QueryRow(`SELECT id, title, description, director, release_year, genre, poster_url, trailer_url FROM movies WHERE id = $1`, id)

	movie := models.Movie{}

	err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Director, &movie.ReleaseYear, &movie.Genre, &movie.PosterURL, &movie.TrailerURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// комментарии для фильма
	rows, err := database.DB.Query(`SELECT c.content, c.created_at, u.username 
		FROM comments c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.movie_id = $1 ORDER BY c.created_at DESC`, id)

	if err != nil {
		http.Error(w, "Ошибка при загрузке комментариев", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		if err := rows.Scan(&comment.Content, &comment.CreatedAt, &comment.Username); err != nil {
			http.Error(w, "Ошибка при загрузке комментариев", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	tmpl := template.Must(template.ParseFiles("web/templates/movieDetail.html"))
	tmpl.Execute(w, struct {
		Movie    models.Movie
		Comments []models.Comment
	}{
		Movie:    movie,
		Comments: comments,
	})
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)
	movieID := mux.Vars(r)["id"]
	content := r.FormValue("content")

	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка получения user_id", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO comments (movie_id, user_id, content) VALUES ($1, $2, $3)", movieID, userID, content)
	if err != nil {
		http.Error(w, "Ошибка при добавлении комментария", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/movie/"+movieID, http.StatusSeeOther)
}
