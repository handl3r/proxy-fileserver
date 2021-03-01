package config

import (
	"os"
	"strconv"
)

func GetInt(key string) int {
	rawValue := os.Getenv(key)
	value, _ := strconv.Atoi(rawValue)
	return value
}

func GetIntWithDefault(key string, defaultValue int) int {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetString(key string) string {
	return os.Getenv(key)
}

func GetStringWitDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
