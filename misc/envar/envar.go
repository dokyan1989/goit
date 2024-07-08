package envar

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}

func GetInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}

	return i
}
