package main

import (
	"log"
	"os"

	"github.com/informeai/cachingo/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error in dotenv: %v", err)
	}
	log.Printf("running in port: %v", os.Getenv("PORT"))
	router := routes.NewRouter()
	log.Fatal(router.Start())
}
