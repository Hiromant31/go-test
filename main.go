package main

import (
	"fmt"
	"go-test/auth"
	"go-test/config"
	"net/http"
)

func main() {
	config.Load()

	http.Handle("/secure", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := auth.GetUserID(r)
		role := auth.GetUserRole(r)
		fmt.Fprintf(w, "Пользователь: %s, роль: %s\n", userID, role)
	})))

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GenerateJWT("user123", "admin")
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(token))
	})

	fmt.Println("Сервер на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
