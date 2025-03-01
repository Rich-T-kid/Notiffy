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
		log.Fatal("loading env file resulted in an error ->", err)
	}
	key := "test"
	fmt.Printf(fmt.Sprintf("enviremtn varible for %s is %s\n", key, os.Getenv(key)))
	fmt.Println("vim-go")
}
