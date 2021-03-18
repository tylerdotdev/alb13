package main

import (
	"log"

	"github.com/alb13/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	server.Start()
}
