package main

import (
	"log"

	"Project/internal/api"
)

func main() {
	log.Println("Application start")

	api.StartServer()

	log.Println("Application terminated")
}
