package handlers

import (
	"html/template"
	"net/http"
)

var users = map[string]string{}

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

	if _, exist := users[username]; exist {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	users[username] = password

	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
