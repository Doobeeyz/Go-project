package handlers

import (
	"html/template"
	"net/http"

	"projectMod/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/registrationPage.html"))
	tmpl.Execute(w, nil)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Введите имя пользователя и пароль", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при шифровании пароля", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
		return
	}

	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
