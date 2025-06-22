package util

import (
	"github.com/joho/godotenv"
)

func LoadEnv(filepath string) {
	err := godotenv.Load(filepath)

	if err != nil {
		Errorf("Error loading .env file, using operating system environment variables.")
	}
}
