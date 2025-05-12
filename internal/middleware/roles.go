package middleware

import (
	"net/http"
	"projectMod/internal/database"
	"projectMod/internal/handlers"
)

func RequireRoleID(requiredRoleID int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := handlers.Store.Get(r, "session-name")
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Error(w, "Войдите в систему", http.StatusUnauthorized)
			return
		}

		username, _ := session.Values["username"].(string)

		var roleID int
		err := database.DB.QueryRow("SELECT role_id FROM users WHERE username = $1", username).Scan(&roleID)
		if err != nil || roleID != requiredRoleID {
			http.Error(w, "Доступ запрещён", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireMinRoleID(minRoleID int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := handlers.Store.Get(r, "session-name")
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Error(w, "Войдите в систему", http.StatusUnauthorized)
			return
		}

		username, _ := session.Values["username"].(string)
		var roleID int
		err := database.DB.QueryRow("SELECT role_id FROM users WHERE username = $1", username).Scan(&roleID)
		if err != nil || roleID < minRoleID {
			http.Error(w, "Доступ запрещён", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
