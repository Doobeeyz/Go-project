package handlers

import (
	"fmt"
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
		SELECT users.id, users.username, roles.name, users.active
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
		Active   bool

	}

	var users []UserInfo
	for rows.Next() {
		var u UserInfo
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Active); err == nil {
			users = append(users, u)
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/roleAdmin.html"))
	tmpl.Execute(w, users)
}

// func BlockUserHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := Store.Get(r, "session-name")
// 	auth, ok := session.Values["authenticated"].(bool)
// 	if !ok || !auth {
// 		http.Error(w, "Войдите в систему", http.StatusUnauthorized)
// 		return
// 	}

// 	username, _ := session.Values["username"].(string)
// 	var user models.User
// 	err := database.DB.QueryRow("SELECT id, username, password, role_id FROM users WHERE username = $1",
// 		username).Scan(&user.ID, &user.Username, &user.Password, &user.RoleID)
// 	if err != nil {
// 		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
// 		return
// 	}

// 	if user.RoleID != 3 {
// 		http.Error(w, "Доступ запрещён", http.StatusForbidden)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
// 		return
// 	}
// 	userID := r.FormValue("user_id")

// 	if userID == "" {
// 		http.Error(w, "ID пользователя не указан", http.StatusBadRequest)
// 		return
// 	}

// 	if userID == fmt.Sprint(user.ID) {
// 		http.Error(w, "Вы не можете заблокировать самого себя", http.StatusForbidden)
// 		return
// 	}

// 	_, err = database.DB.Exec("UPDATE users SET active = false WHERE id = $1", userID)
// 	if err != nil {
// 		http.Error(w, "Ошибка при блокировке пользователя", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/superadmin", http.StatusSeeOther)
// }

func ToggleUserActiveHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Войдите в систему", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)
	var currentUser models.User
	err := database.DB.QueryRow("SELECT id, role_id FROM users WHERE username = $1", username).
		Scan(&currentUser.ID, &currentUser.RoleID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	if currentUser.RoleID != 3 {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("user_id")
	active := r.FormValue("active") 

	if userID == "" || active == "" {
		http.Error(w, "Данные не указаны", http.StatusBadRequest)
		return
	}

	if userID == fmt.Sprint(currentUser.ID) {
		http.Error(w, "Вы не можете изменить свой статус", http.StatusForbidden)
		return
	}

	newStatus := (active != "true")
	_, err = database.DB.Exec("UPDATE users SET active = $1 WHERE id = $2", newStatus, userID)
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
}


