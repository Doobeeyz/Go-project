package handlers

import(
	"net/http"
	"html/template"
)



func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "Доступ запрещен. Войдите в систему.", http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/mainPage.html"))
	tmpl.Execute(w, nil)
}
