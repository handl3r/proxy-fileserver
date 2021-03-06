package config

import (
	"os"
	"proxy-fileserver/enums"
	"strconv"
	"time"
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

func GetTimeDuration(key string) (time.Duration, error) {
	rawValue := os.Getenv(key)
	if len(rawValue) < 2 {
		return 0, enums.ErrInValidConfigTimeDuration
	}
	sizeStr := rawValue[0 : len(rawValue)-1]
	unit := string(rawValue[len(rawValue)-1])
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, enums.ErrInValidConfigTimeDuration
	}
	var timeDuration time.Duration
	switch unit {
	case "H":
		timeDuration = time.Duration(size) * time.Hour
	case "M":
		timeDuration = time.Duration(size) * time.Minute
	case "S":
		timeDuration = time.Duration(size) * time.Second
	default:
		err = enums.ErrInValidConfigTimeDuration
	}
	if err != nil {
		return 0, err
	}
	return timeDuration, nil
}
