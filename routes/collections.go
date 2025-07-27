package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go-test/appwrite"
	"go-test/auth"
)

type createCollectionRequest struct {
	Name string `json:"name"`
}

func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	role := auth.GetUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden: недостаточно прав", http.StatusForbidden)
		return
	}

	var req createCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
		return
	}
	jwtToken := auth.GetRawJWT(r)
	client := appwrite.NewClientWithJWT(jwtToken)

	collection, err := client.CreateCollection(
		req.Name,
		[]string{},
		false,
		true,
	)
	if err != nil {
		log.Println("CreateCollection error:", err)
		http.Error(w, "Ошибка Appwrite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)

}
