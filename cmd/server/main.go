package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"projectMod/internal/database"
	"projectMod/internal/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	wd, _ := os.Getwd()
	fs := http.FileServer(http.Dir(filepath.Join(wd, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	database.Init()

	r := routes.NewRouter()

	fmt.Println("server is running on port: 8080")

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_PORT"), r))
}
