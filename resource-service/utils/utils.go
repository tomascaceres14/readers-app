package utils

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	fmt.Println("found " + key)
	return fallback
}
