package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env pada root folder
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error/Tidak menemukan file .env, menggunakan environment variable bawaan")
	}
}
