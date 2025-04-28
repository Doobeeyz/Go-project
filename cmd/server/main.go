package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projectMod/internal/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	r := routes.NewRouter()

	fmt.Println("server is running on port: 8080")

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), r))
}
