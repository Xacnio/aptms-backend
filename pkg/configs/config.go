package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

func Gets(keys ...*string) {
	for i := range keys {
		*keys[i] = os.Getenv(*keys[i])
	}
}
