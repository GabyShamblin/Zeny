package main

import (
	//"fmt"
	//"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := "8080"
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Print("Server starting at localhost:8080")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
