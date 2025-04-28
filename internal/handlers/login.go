package handlers

import(
	"net/http"
	"html/template"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func LoginPage(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("web/templates/loginPage.html"))
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = true
	session.Values["username"] = "user1"
	session.Save(r, w)

    http.Redirect(w, r, "/main", http.StatusSeeOther) 
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}