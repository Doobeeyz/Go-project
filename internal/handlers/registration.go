package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"projectMod/internal/database"

	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	letterRegex   = regexp.MustCompile(`[A-Za-z]`)
	digitRegex    = regexp.MustCompile(`[0-9]`)
	specialRegex  = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\:;"',.?/~]`)
	disallowed    = regexp.MustCompile(`[<>]`)
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

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Введите имя пользователя и пароль", http.StatusBadRequest)
		return
	}

	if !usernameRegex.MatchString(username) {
		http.Error(w, "Имя пользователя должно содержать от 3 до 20 символов (латинские буквы, цифры, подчёркивание)", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Пароль должен содержать не менее 6 символов", http.StatusBadRequest)
		return
	}

	if disallowed.MatchString(password) {
		http.Error(w, "Пароль не должен содержать символы < или >", http.StatusBadRequest)
		return
	}

	if !letterRegex.MatchString(password) {
		http.Error(w, "Пароль должен содержать хотя бы одну латинскую букву", http.StatusBadRequest)
		return
	}

	if !digitRegex.MatchString(password) {
		http.Error(w, "Пароль должен содержать хотя бы одну цифру", http.StatusBadRequest)
		return
	}

	if !specialRegex.MatchString(password) {
		http.Error(w, "Пароль должен содержать хотя бы один специальный символ", http.StatusBadRequest)
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
