package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")

	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Доступ запрещен. Войдите в систему.", http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/mainPage.html"))
	tmpl.Execute(w, nil)
}
