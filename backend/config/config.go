package config

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// EnvVariable - load configs from .env file
func EnvVariable(key string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..")

	// load .env file
	err := godotenv.Load(basepath + "/.env")

	if err != nil {
		log.Println("Cannot load .env file, defaulting to env variables")
	}

	return os.Getenv(key)
}
