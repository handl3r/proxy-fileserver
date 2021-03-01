package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvironments() {
	err := godotenv.Load()
	if os.IsNotExist(err) {
		panic("File .env is missing")
	} else if err != nil {
		panic(fmt.Sprintf("Error when load file .env: [%v]", err))
	}
}
