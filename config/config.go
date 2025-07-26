package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppwriteHost    string
	AppwriteKey     string
	AppwriteProject string
	JWTSecret       string
}

var Cfg *Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, читаю переменные из окружения")
	}

	Cfg = &Config{
		AppwriteHost:    os.Getenv("APPWRITE_HOST"),
		AppwriteKey:     os.Getenv("APPWRITE_KEY"),
		AppwriteProject: os.Getenv("APPWRITE_PROJECT_ID"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
	}

	if Cfg.AppwriteHost == "" || Cfg.AppwriteKey == "" || Cfg.JWTSecret == "" || Cfg.AppwriteProject == "" {
		log.Fatal("Не заданы все обязательные переменные окружения")
	}
}
