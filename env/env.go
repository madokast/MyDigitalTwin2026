package env

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		slog.Warn("failed to load env", "file", file, "err", err)
	}
	return err
}

func CheckEnv(names []string) error {
	for _, name := range names {
		if os.Getenv(name) == "" {
			return fmt.Errorf("env %s miss", name)
		}
	}
	return nil
}
