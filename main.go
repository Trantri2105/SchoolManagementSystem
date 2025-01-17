package main

import (
	"SchoolManagement/repo"
	"SchoolManagement/transport"
	"github.com/joho/godotenv"
	"log"
)

func loadEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnvVariable()
	db := repo.PostgresConnect()
	redisClient := repo.RedisConnect()

	r := transport.NewHttpServer(db, redisClient)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
