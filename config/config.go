package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}
	FetchURL string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config := Config{}
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")
	config.FetchURL = os.Getenv("FETCH_URL")
	return &config, nil
}
