package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error when loading file configuration: ", err.Error())
	}

	expInt, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host:   os.Getenv("DB_HOST"),
			Port:   os.Getenv("DB_PORT"),
			User:   os.Getenv("DB_USER"),
			Pass:   os.Getenv("DB_PASS"),
			Name:   os.Getenv("DB_NAME"),
			Schema: os.Getenv("DB_SCHEMA"),
			Tz:     os.Getenv("DB_TZ"),
		},
		Jwt: Jwt{
			Key: os.Getenv("JWT_KEY"),
			Exp: expInt,
		},
		FileUpload: FileUpload{
			Path: os.Getenv("FILE_PATH"),
		},
		Redis: Redis{
			Address:  os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASS"),
		},
	}
}
