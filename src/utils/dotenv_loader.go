package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}
