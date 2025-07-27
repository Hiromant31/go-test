package auth

import (
	"encoding/json"
	"fmt"
	"go-test/appwrite"
	"log"
	"net/http"
	"sync"

	"github.com/appwrite/sdk-for-go/id"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	JWT   string `json:"jwt,omitempty"`
	Error string `json:"error,omitempty"`
}

var mu sync.Mutex
var roles = map[string]string{} // userId -> роль

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	client := appwrite.NewAppwriteClient()
	userId := id.Unique()
	_, err := client.CreateUser(userId, req.Email, req.Password)
	if err != nil {
		writeJSON(w, authResponse{Error: "Ошибка создания пользователя: " + err.Error()})
		return
	}
	fmt.Println("Creating user with userId:", userId, "length:", len(userId))

	session, err := client.LoginUser(req.Email, req.Password)
	if err != nil {
		writeJSON(w, authResponse{Error: "Ошибка входа: " + err.Error()})
		return
	}

	// Генерация JWT
	jwt, err := client.CreateJWT()
	if err != nil {
		writeJSON(w, authResponse{Error: "Ошибка генерации JWT: " + err.Error()})
		return
	}

	mu.Lock()
	roles[session.UserId] = "user"
	mu.Unlock()

	// Устанавливаем JWT в заголовке ответа
	w.Header().Set("Authorization", "Bearer "+jwt.Jwt)
	writeJSON(w, authResponse{JWT: jwt.Jwt})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	client := appwrite.NewAppwriteClient()
	session, err := client.LoginUser(req.Email, req.Password)
	fmt.Println(req.Email, req.Password)
	if err != nil {
		log.Printf("Ошибка при входе пользователя %s: %v", req.Email, err)
		writeJSON(w, authResponse{Error: "Неверный логин или пароль"})
		return
	}

	// Генерация JWT
	jwt, err := client.CreateJWT()
	if err != nil {
		log.Printf("Ошибка при генерации JWT для пользователя %s: %v", req.Email, err)
		writeJSON(w, authResponse{Error: "Ошибка генерации токена"})
		return
	}

	mu.Lock()
	role := roles[session.UserId]
	mu.Unlock()
	fmt.Println("Role:", role)

	// Устанавливаем JWT в заголовке ответа
	w.Header().Set("Authorization", "Bearer "+jwt.Jwt)
	writeJSON(w, authResponse{JWT: jwt.Jwt})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
