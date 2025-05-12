package handlers

import (
	"net/http"
	"projectMod/internal/database"
	"strconv"
)

func ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Войдите в систему", http.StatusUnauthorized)
		return
	}

	username, _ := session.Values["username"].(string)

	var roleID int
	err := database.DB.QueryRow("SELECT role_id FROM users WHERE username = $1", username).Scan(&roleID)
	if err != nil {
		http.Error(w, "Ошибка получения роли", http.StatusInternalServerError)
		return
	}

	if roleID != 3 {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	// Чтение значений из формы
	targetUserID := r.FormValue("user_id")
	newRoleIDStr := r.FormValue("role_id")

	newRoleID, err := strconv.Atoi(newRoleIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID роли", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET role_id = $1 WHERE id = $2", newRoleID, targetUserID)
	if err != nil {
		http.Error(w, "Ошибка обновления роли", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
}
