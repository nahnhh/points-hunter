package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	email := os.Getenv("EMAIL")
	pass := os.Getenv("PASS")
	token := os.Getenv("GEMINI_API_KEY")

	fmt.Printf("%v %v", email, pass)
	GetEvent(token)
}
