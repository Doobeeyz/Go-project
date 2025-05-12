package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/models"
)

func SuperAdminRoleHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Войдите в систему", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, password, role_id FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.RoleID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	if user.RoleID != 3 {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	rows, err := database.DB.Query(`
        SELECT users.id, users.username, roles.name 
        FROM users 
        JOIN roles ON users.role_id = roles.id`)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserInfo struct {
		ID       int
		Username string
		Role     string
	}

	var users []UserInfo
	for rows.Next() {
		var u UserInfo
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err == nil {
			users = append(users, u)
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/roleAdmin.html"))
	tmpl.Execute(w, users)
}
