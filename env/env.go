package env

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	return godotenv.Load(file)
}

func CheckEnv(names []string) error {
	for _, name := range names {
		if os.Getenv(name) == "" {
			return fmt.Errorf("env %s miss", name)
		}
	}
	return nil
}

func Get(name string) (string, bool) {
	value := os.Getenv(name)
	if value == "" {
		return "", false
	}
	return value, true
}

func MustGet(name string) string {
	value, ok := Get(name)
	if !ok {
		slog.Error("env not set", "name", name)
		os.Exit(1)
	}
	return value
}
