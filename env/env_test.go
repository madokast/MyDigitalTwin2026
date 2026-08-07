package env

import (
	"testing"
)

func TestLoadEnv(t *testing.T) {
	err := LoadEnv("testdata/.env")
	if err != nil {
		t.Error(err)
	}
}

func TestCheckEnv(t *testing.T) {
	err := LoadEnv("testdata/.env")
	if err != nil {
		t.Error(err)
	}
	err = CheckEnv([]string{"DT_ENV"})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckEnvMiss(t *testing.T) {
	err := LoadEnv("testdata/.env")
	if err != nil {
		t.Error(err)
	}
	err = CheckEnv([]string{"DT_ENV_2"})
	if err == nil {
		t.Error("should miss env")
	}
}
