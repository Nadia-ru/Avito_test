package utilities

import "github.com/joho/godotenv"

func CheckEnvFile() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

}
