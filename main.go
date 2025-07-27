package main

import (
	"encoding/json"
	"fmt"
	"go-test/appwrite"
	"go-test/auth"
	"go-test/config"
	"go-test/routes"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	config.Load()

	mux := http.NewServeMux()

	mux.Handle("/secure", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.GetUserID(r)
		role := auth.GetUserRole(r)
		resp := map[string]string{
			"message": fmt.Sprintf("Пользователь: %s, роль: %s", userID, role),
		}
		json.NewEncoder(w).Encode(resp)
	})))

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GenerateJWT("user123", "admin")
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(token))
	})

	mux.Handle("/collections", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := auth.GetRawJWT(r)
		client := appwrite.NewClientWithJWT(jwtToken)
		data, err := client.GetCollections(nil, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	})))

	mux.HandleFunc("/api/register", auth.RegisterHandler)

	mux.HandleFunc("/api/login", auth.LoginHandler)

	mux.Handle("/collections/create", auth.JWTMiddleware(http.HandlerFunc(routes.CreateCollectionHandler)))

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // фронтенд-адрес
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	fmt.Println("Сервер на http://localhost:8080")
	http.ListenAndServe(":8080", corsHandler)
}
