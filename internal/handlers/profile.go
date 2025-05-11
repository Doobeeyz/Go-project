package handlers

import (
	"html/template"
	"net/http"
	"projectMod/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")

	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username, _ := session.Values["username"].(string)

	flash, _ := session.Values["flash"].(string)
	delete(session.Values, "flash")
	session.Save(r, w)

	tmpl := template.Must(template.ParseFiles("web/templates/profilePage.html"))
	tmpl.Execute(w, struct {
		Username string
		Flash string
	}{
		Username: username,
		Flash:    flash,
	})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}
	username, _ := session.Values["username"].(string)

	oldPassword := r.FormValue("oldPassword")
	newPassword := r.FormValue("newPassword")

	var storedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(oldPassword))
	if err != nil {
		session.Values["flash"] = "Неверный старый пароль"
		session.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		session.Values["flash"] = "Ошибка хеширования пароля"
		session.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET password = $1 WHERE username = $2", string(hashedNewPassword), username)
	if err != nil {
		session.Values["flash"] = "Не удалось обновить пароль"
		session.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func ChangeUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldUsername, _ := session.Values["username"].(string)
	newUsername := r.FormValue("newUsername")

	if newUsername == "" {
		session.Values["flash"] = "Имя пользователя не может быть пустым"
		session.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	_, err := database.DB.Exec("UPDATE users SET username = $1 WHERE username = $2", newUsername, oldUsername)
	if err != nil {
		session.Values["flash"] = "Ошибка при обновлении имени пользователя"
		session.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	session.Values["username"] = newUsername
	session.Values["flash"] = "Имя пользователя успешно изменено"
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
